package main

import (
	"testing"
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
