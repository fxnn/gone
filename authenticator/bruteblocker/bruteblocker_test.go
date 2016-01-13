package bruteblocker

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGrowingDelay(t *testing.T) {

	var max = 10 * time.Second
	var step = 1 * time.Second
	var sut = New(max, step)
	var results = make([]time.Duration, 0)

	for i := 0; i < 3; i++ {
		results = append(results, sut.Delay("user", "ip", false))
	}

	if !reflect.DeepEqual(results, seconds(0, 1, 2)) {
		t.Fatalf("Expected Delay results to grow by 1 second starting at 0, but was %v", results)
	}

}

func TestDelayLimitedAtMax(t *testing.T) {

	var max = 10 * time.Second
	var step = 1 * time.Second
	var sut = New(max, step)
	var results = make([]time.Duration, 0)

	for i := 0; i <= 20; i++ {
		results = append(results, sut.Delay("user", "ip", false))
	}

	if results[20] != max {
		t.Fatalf("Expected 21th delay to be maximum, but was %v", results[20])
	}

}

func TestDelayForUserGrowsAsSpecified(t *testing.T) {

	var max = 10 * time.Second
	var step = 1 * time.Second
	var sut = New(max, step)
	var results = make([]time.Duration, 0)

	for i := 0; i <= 10; i++ {
		results = append(results, sut.Delay("user", fmt.Sprintf("ip%d", i), false))
	}

	if results[10] != max {
		t.Fatalf("Expected 11th delay to be maximum, but was %v", results[10])
	}

}

func TestDelayForAddrGrowsSlower(t *testing.T) {

	var max = 10 * time.Second
	var step = 1 * time.Second
	var sut = New(max, step)
	var results = make([]time.Duration, 0)

	for i := 0; i <= 10; i++ {
		results = append(results, sut.Delay(fmt.Sprintf("user%d", i), "ip", false))
	}

	if results[10] >= max {
		t.Fatalf("Expected 11th delay to be less than maximum, but was %v", results[10])
	}
	if results[10] != 1*time.Second {
		t.Fatalf("Expected 11th delay to be one second, but was %v", results[10])
	}

}

func TestDelayForGlobalGrowsSlower(t *testing.T) {

	var max = 10 * time.Second
	var step = 1 * time.Second
	var sut = New(max, step)
	var results = make([]time.Duration, 0)

	for i := 0; i <= 10; i++ {
		results = append(results, sut.Delay(fmt.Sprintf("user%d", i), fmt.Sprintf("ip%d", i), false))
	}

	if results[10] >= max {
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
