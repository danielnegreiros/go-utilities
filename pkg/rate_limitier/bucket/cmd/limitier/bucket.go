package main

import (
	"errors"
	"sync"
	"time"
)

const (
	errMsgInputValidation = "refill rate and bucket size must be greater than 0"
)

type bucketLimiter struct {
	refillRateMinute int
	size             int
	buckets          map[string]*bucket
}

type bucket struct {
	key          string
	currentUnits int
	lastSeen     time.Time
	mutex        sync.Mutex
}

func NewBucketLimiter(refillRateMin int, size int) (*bucketLimiter, error) {

	if refillRateMin <= 0 || size <= 0 {
		return nil, errors.New(errMsgInputValidation)
	}

	return &bucketLimiter{
		refillRateMinute: refillRateMin,
		size:             size,
		buckets:          make(map[string]*bucket),
	}, nil
}
