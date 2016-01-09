package bruteblocker

import "time"

type request struct {
	userId     string
	sourceAddr string
	successful bool
	response   chan<- time.Duration
}

type BruteBlocker struct {
	requests            chan request
	shutdown            chan struct{}
	countFailedAttempts map[string]int
	lastFailedAttempt   map[string]time.Time
	delayMax            time.Duration
	delayStep           time.Duration
	globalDelayStep     time.Duration
}

// New creates a new BruteBlocker instance.
// This starts a goroutine, which has to be shut down eventually!
//
// delayMax denotes the maximum delay to impose, delayStep denotes the delay
// increment per failed authentication attempt.
func New(delayMax time.Duration, delayStep time.Duration) *BruteBlocker {
	var result = &BruteBlocker{
		requests:            make(chan request),
		shutdown:            make(chan struct{}),
		countFailedAttempts: make(map[string]int),
		lastFailedAttempt:   make(map[string]time.Time),
		delayMax:            delayMax,
		delayStep:           delayStep,
		// NOTE, that the global delay grows slower
		globalDelayStep: delayStep / 10,
	}
	go result.serve()
	return result
}

func (b *BruteBlocker) serve() {
	for {
		select {
		case rq := <-b.requests:
			rq.response <- b.delay(rq.userId, rq.sourceAddr, rq.successful)
		case _ = <-b.shutdown:
			break
		}
	}
}

func (b *BruteBlocker) delay(userId string, sourceAddr string, successful bool) time.Duration {
	// NOTE, that we also impose a delay on successful authentication attempts,
	// so that the attacker needs our response.

	var userDelay = b.delayFor(userId, b.delayStep, successful)
	var addrDelay = b.delayFor(sourceAddr, b.delayStep, successful)
	var globalDelay = b.delayFor("", b.globalDelayStep, successful)

	var maxDelay = max(globalDelay, max(userDelay, addrDelay))
	return min(maxDelay, b.delayMax)
}

func (b *BruteBlocker) delayFor(id string, step time.Duration, successful bool) time.Duration {
	var count = b.countFailedAttempts[id]
	if !successful {
		b.countFailedAttempts[id] = count + 1
		b.lastFailedAttempt[id] = time.Now()
	}

	return step * time.Duration(count)
}

// max returns maximum between two values.
// cf. http://stackoverflow.com/a/28207228/3281722 why this isn't done with a
// stdlib call (tl;dr: there is none).
func max(a time.Duration, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

// min returns minimum between two values.
// cf. http://stackoverflow.com/a/28207228/3281722 why this isn't done with a
// stdlib call (tl;dr: there is none).
func min(a time.Duration, b time.Duration) time.Duration {
	if a > b {
		return b
	}
	return a
}

func (b *BruteBlocker) ShutDown() {
	b.shutdown <- struct{}{}
}

func (b *BruteBlocker) Delay(userId string, sourceAddr string, successful bool) time.Duration {
	var response = make(chan time.Duration)
	b.requests <- request{userId, sourceAddr, successful, response}
	return <-response
}
