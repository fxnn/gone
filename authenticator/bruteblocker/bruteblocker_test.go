package bruteblocker

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var (
	defaultMax        = 10 * time.Second
	defaultUserStep   = 1 * time.Second
	defaultAddrStep   = 100 * time.Millisecond
	defaultGlobalStep = 50 * time.Millisecond
	defaultDropAfter  = 10 * time.Second
)

func newSut() *BruteBlocker {
	return newSutDroppingAfter(defaultDropAfter)
}

func newSutDroppingAfter(dropAfter time.Duration) *BruteBlocker {
	return New(
		defaultMax,
		defaultUserStep,
		defaultAddrStep,
		defaultGlobalStep,
		dropAfter)
}

func TestCleanUp(t *testing.T) {

	var sut = newSutDroppingAfter(0)

	sut.Delay("user", "ip", false)
	time.Sleep(100 * time.Millisecond)
	sut.CleanUp()

	if delay := sut.Delay("user", "ip", false); delay != 0 {
		t.Fatalf("Expected delay after cleanup to be 0, but was %v", delay)
	}

}

func TestDelayPanicsAfterShutdown(t *testing.T) {

	var sut = newSut()
	sut.ShutDown()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Panic expected; but actually the test ended regularly")
		}
	}()

	sut.Delay("", "", false)

}

func TestShutdownPanicsAfterShutdown(t *testing.T) {

	var sut = newSut()
	sut.ShutDown()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Panic expected; but actually the test ended regularly")
		}
	}()

	sut.ShutDown()

}

func TestGrowingDelay(t *testing.T) {

	var sut = newSut()
	var results = make([]time.Duration, 0)

	for i := 0; i < 3; i++ {
		results = append(results, sut.Delay("user", "ip", false))
	}

	if !reflect.DeepEqual(results, seconds(0, 1, 2)) {
		t.Fatalf("Expected Delay results to grow by 1 second starting at 0, but was %v", results)
	}

}

func TestDelayLimitedAtMax(t *testing.T) {

	var sut = newSut()
	var results = make([]time.Duration, 0)

	for i := 0; i <= 20; i++ {
		results = append(results, sut.Delay("user", "ip", false))
	}

	if results[20] != defaultMax {
		t.Fatalf("Expected 21th delay to be maximum, but was %v", results[20])
	}

}

func TestDelayForUserGrowsAsSpecified(t *testing.T) {

	var sut = newSut()
	var results = make([]time.Duration, 0)

	for i := 0; i <= 10; i++ {
		results = append(results, sut.Delay("user", fmt.Sprintf("ip%d", i), false))
	}

	if results[10] != defaultMax {
		t.Fatalf("Expected 11th delay to be maximum, but was %v", results[10])
	}

}

func TestDelayForAddrGrowsSlower(t *testing.T) {

	var sut = newSut()
	var results = make([]time.Duration, 0)

	for i := 0; i <= 10; i++ {
		results = append(results, sut.Delay(fmt.Sprintf("user%d", i), "ip", false))
	}

	if results[10] >= defaultMax {
		t.Fatalf("Expected 11th delay to be less than maximum, but was %v", results[10])
	}
	if results[10] != 1*time.Second {
		t.Fatalf("Expected 11th delay to be one second, but was %v", results[10])
	}

}

func TestDelayForGlobalGrowsSlower(t *testing.T) {

	var sut = newSut()
	var results = make([]time.Duration, 0)

	for i := 0; i <= 10; i++ {
		results = append(results, sut.Delay(fmt.Sprintf("user%d", i), fmt.Sprintf("ip%d", i), false))
	}

	if results[10] >= defaultMax {
		t.Fatalf("Expected 11th delay to be less than maximum, but was %v", results[10])
	}
	if results[10] != 500*time.Millisecond {
		t.Fatalf("Expected 11th delay to be 500 milliseconds, but was %v", results[10])
	}

}

func seconds(values ...time.Duration) []time.Duration {
	var result = make([]time.Duration, len(values))
	for i, v := range values {
		result[i] = v * time.Second
	}
	return result
}
