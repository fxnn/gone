package bruteblocker

import "time"

type request struct {
	userId     string
	sourceAddr string
	successful bool
	response   chan<- time.Duration
}

// BruteBlocker encapsulates data and behaviour for brute force attack
// detection.
// It memorizes how many failed attempts occured per user and ip address, and
// when the last one occured.
// To allow for multi threaded access, it also creates a goroutine and stores
// channels for messaging.
type BruteBlocker struct {
	requests            chan request
	shutdown            chan struct{}
	countFailedAttempts map[string]int
	lastFailedAttempt   map[string]time.Time
	delayMax            time.Duration
	userDelayStep       time.Duration
	addrDelayStep       time.Duration
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
		userDelayStep:       delayStep,
		addrDelayStep:       max(1, delayStep/10),
		globalDelayStep:     max(1, delayStep/20),
	}
	go result.serve()
	return result
}

// ShutDown stops the goroutine associated with this BruteBlocker instance.
// The instance is no longer functional.
func (b *BruteBlocker) ShutDown() {
	b.shutdown <- struct{}{}
}

// Delay informs this BruteBlocker about a login attempt and returns the
// amount of time the user should be blocked.
//
// Note that, when delaying the response while allowing concurrent requests, you
// should also delay after an successful authentication.
// This way, the attacker needs to await the full delay in order to know whether
// his authentication attempt succeeded or not.
func (b *BruteBlocker) Delay(userId string, sourceAddr string, successful bool) time.Duration {
	var response = make(chan time.Duration)
	b.requests <- request{userId, sourceAddr, successful, response}
	return <-response
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
	var userDelay = b.delayCriterion("user="+userId, b.userDelayStep, successful)
	var addrDelay = b.delayCriterion("addr="+sourceAddr, b.addrDelayStep, successful)
	var globalDelay = b.delayCriterion("global", b.globalDelayStep, successful)

	var maxDelay = max(globalDelay, max(userDelay, addrDelay))
	return min(maxDelay, b.delayMax)
}

func (b *BruteBlocker) delayCriterion(id string, step time.Duration, successful bool) time.Duration {
	// NOTE, that we also impose a delay on successful authentication attempts,
	// so that the attacker needs our response.

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
