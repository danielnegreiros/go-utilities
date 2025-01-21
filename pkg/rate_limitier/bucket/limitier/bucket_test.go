package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestNewBucketLimiter(t *testing.T) {
	// Test valid input
	limiter, err := NewBucketLimiter(10, 100)
	if err != nil {
		t.Error(err)
	}
	if limiter.refillRateMinute != 10 {
		t.Errorf("expected RefillRateMinute = 10, got %d", limiter.refillRateMinute)
	}
	if limiter.size != 100 {
		t.Errorf("expected Size = 100, got %d", limiter.size)
	}
	if len(limiter.buckets) != 0 {
		t.Errorf("expected Buckets to be empty, got %d elements", len(limiter.buckets))
	}

	failingScenarrions := [][]int{
		{0, 1},
		{1, 0},
		{-1, 1},
		{1, -1},
	}

	for _, v := range failingScenarrions {
		_, err = NewBucketLimiter(v[0], v[1])
		if err == nil {
			t.Error("should have failed")
		}
	}
}

func TestRequestNewBucket(t *testing.T) {
	bl, err := NewBucketLimiter(60, 80)
	if err != nil {
		t.Error(err)
	}

	isAllowed := bl.RequestToken("keya")
	if !isAllowed {
		t.Error("expected to be allowed")
	}
}

func TestRequestNotAllowd(t *testing.T) {
	bl, err := NewBucketLimiter(1, 1)
	if err != nil {
		t.Error(err)
	}

	isAllowed := bl.RequestToken("keya")
	if !isAllowed {
		t.Error("expected to be allowd")
	}

	isAllowed = bl.RequestToken("keya")
	if isAllowed {
		t.Error("expected to be denied")
	}

}

func TestRefillReachTop(t *testing.T) {

	top := 11
	refillRate := 60 // 1 unit per second

	b := bucket{
		key:          "a",
		lastRefilled: time.Now(),
		currentUnits: 1,
		mutex:        sync.Mutex{},
	}

	b.refill(time.Now().Add(15*time.Second), top, refillRate)
	if b.currentUnits != top {
		t.Errorf("\nExpected: %d Found: %d", top, b.currentUnits)
	}

}

func TestPartialRefill(t *testing.T) {
	top := 1000
	refillRate := 120 // 2 unit per second
	currUnit := 10
	elapsedSeconds := 15
	expected := currUnit + (2 * elapsedSeconds)

	b := bucket{
		key:          "a",
		lastRefilled: time.Now(),
		currentUnits: currUnit,
		mutex:        sync.Mutex{},
	}

	b.refill(time.Now().Add(time.Duration(elapsedSeconds)*time.Second), top, refillRate)
	if b.currentUnits != expected {
		t.Errorf("\nExpected: %d Found: %d", expected, b.currentUnits)
	}
}

func TestSmallReffil(t *testing.T) {
	top := 10
	refillRate := 30 // 0.5 unit per second
	currUnit := 0
	expected := 3

	b := bucket{
		key:          "a",
		lastRefilled: time.Now(),
		currentUnits: currUnit,
		mutex:        sync.Mutex{},
	}

	b.refill(time.Now().Add(6*time.Second), top, refillRate)
	if b.currentUnits != expected {
		t.Errorf("\nExpected: %d Found: %d", expected, b.currentUnits)
	}
}
