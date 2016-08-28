package bruteblocker

import "time"

// BruteBlocker encapsulates data and behaviour for brute force attack
// detection.
// It memorizes how many failed attempts occured per user and ip address, and
// when the last one occured.
// To allow for multi threaded access, it also creates a goroutine and stores
// channels for messaging.
type BruteBlocker struct {
	// channels
	requests      chan request
	cleanUpTicker *time.Ticker

	// state
	shutdown            bool
	countFailedAttempts map[string]int
	lastFailedAttempt   map[string]time.Time

	// configuration
	delayMax        time.Duration
	userDelayStep   time.Duration
	addrDelayStep   time.Duration
	globalDelayStep time.Duration
	dropAfter       time.Duration
}

type request func()

// New creates a new BruteBlocker instance.
// This starts a goroutine, which might be shut down at any time using the
// ShutDown() func.
//
// delayMax denotes the maximum delay to impose, and the delayStep values denote
// the delay increment per failed authentication attempt.
// Separate delays are tracked per user, per source ip address and globally (i.e.
// for all users and ip addresses).
// dropAfter denotes how long after the last failed login attempt the delay
// should be dropped.
func New(
	delayMax time.Duration,
	userDelayStep time.Duration,
	addrDelayStep time.Duration,
	globalDelayStep time.Duration,
	dropAfter time.Duration,
) *BruteBlocker {
	var cleanUpInterval = max(1*time.Second, min(60*time.Second, dropAfter))
	var result = &BruteBlocker{
		requests:            make(chan request),
		cleanUpTicker:       time.NewTicker(cleanUpInterval),
		shutdown:            false,
		countFailedAttempts: make(map[string]int),
		lastFailedAttempt:   make(map[string]time.Time),
		delayMax:            delayMax,
		userDelayStep:       userDelayStep,
		addrDelayStep:       addrDelayStep,
		globalDelayStep:     globalDelayStep,
		dropAfter:           dropAfter,
	}
	go result.serve()
	go result.cleanUpEachTick()
	return result
}

// ShutDown stops the goroutine associated with this BruteBlocker instance.
// The instance is no longer functional and will panic on use.
func (b *BruteBlocker) ShutDown() {
	b.requests <- func() {
		b.shutdown = true
	}
}

// CleanUp removes old entries from memory.
// An entry is old if the last failed login attempt is older than the dropAfter
// parameter requires.
func (b *BruteBlocker) CleanUp() {
	var response = make(chan struct{})
	for id, timestamp := range b.lastFailedAttempt {
		if time.Since(timestamp) > b.dropAfter {
			// NOTE that delete must happen as request, just like all accesses
			// to the maps.
			// HINT: We synchronize calls using response to not surprise the
			// with not being done with cleanup after CleanUp() returns
			b.requests <- func() {
				delete(b.lastFailedAttempt, id)
				delete(b.countFailedAttempts, id)
				response <- struct{}{}
			}
			<-response
		}
	}
}

// Delay informs this BruteBlocker about a login attempt and returns the
// amount of time the user should be blocked.
//
// Note that, when delaying the response while allowing concurrent requests, you
// should also delay after an successful authentication.
// This way, the attacker needs to await the full delay in order to know whether
// his authentication attempt succeeded or not.
func (b *BruteBlocker) Delay(userID string, sourceAddr string, successful bool) time.Duration {
	var response = make(chan time.Duration)
	// NOTE that delay call must happen as request, just as all accesses to
	// internal data structures
	b.requests <- func() {
		response <- b.delay(userID, sourceAddr, successful)
	}
	return <-response
}

// serve accepts and runs requests from the BruteBlocker struct. This way, all
// accesses to the not thread-safe maps can happen in one single gorotuine
// (communication by sharing).
func (b *BruteBlocker) serve() {
	for {
		select {
		case rq := <-b.requests:
			rq()
			if b.shutdown {
				close(b.requests)
				b.cleanUpTicker.Stop()
				return
			}
		}
	}
}

func (b *BruteBlocker) cleanUpEachTick() {
	for range b.cleanUpTicker.C {
		b.CleanUp()
	}
}

// delay calculates the delay for the given login attempt and updates internal
// data structures.
// It MUST be called as request.
func (b *BruteBlocker) delay(userID string, sourceAddr string, successful bool) time.Duration {
	var userDelay = b.delayCriterion("user="+userID, b.userDelayStep, successful)
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
