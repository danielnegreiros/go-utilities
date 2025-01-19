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

// refill increses current units if they are available
// this operation happens on a request level
// avoiding long running term routine
func (b *bucket) refill(nowTime time.Time, top int, rate int) error {
	elaspsed := int(nowTime.Sub(b.lastSeen).Seconds())
	secRate := rate / 60

	refill := elaspsed * secRate

	if b.currentUnits + refill >= top {
		b.currentUnits = top
		return nil
	}

	b.currentUnits += refill
	return nil
}

func (b *bucket) cashOut() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.currentUnits--
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

// RequestToken is the entry point for the bucket limitier
func (bl *bucketLimiter) RequestToken(key string) bool {

	bucket, ok := bl.buckets[key]
	if !ok {
		bucket = bl.createBucket(key)
		bl.buckets[key] = bucket
	}

	bucket.refill(time.Now(), bl.size, bl.refillRateMinute)

	if bucket.currentUnits == 0 {
		return false
	}

	bucket.cashOut()
	return true
}

// createBucket create a new bucket for a key if it doest exist
func (bl *bucketLimiter) createBucket(key string) *bucket {

	return &bucket{
		key:          key,
		currentUnits: bl.size,
		lastSeen:     time.Now(),
		mutex:        sync.Mutex{},
	}
}
