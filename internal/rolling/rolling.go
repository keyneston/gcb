package rolling

import (
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type timeGetter func() time.Time

type bucket struct {
	attempts uint
	failures uint
}

func (b *bucket) clear() {
	b.attempts = 0
	b.failures = 0
}

type rolling struct {
	Buckets     []*bucket
	StartBucket int
	StartTime   time.Time
	Mutex       *sync.RWMutex
	count       int

	time timeGetter
}

func New(count int) *rolling {
	buckets := make([]*bucket, count)
	for i := range buckets {
		buckets[i] = &bucket{}
	}

	return &rolling{
		Buckets: buckets,
		Mutex:   &sync.RWMutex{},

		count: count,

		// time.Now is used except for testing
		time: time.Now,
	}
}

func (r *rolling) current() *bucket {
	now := r.time()
	diff := int(now.Sub(r.StartTime) / time.Second)

	// If its been more than bCount seconds since
	if diff > r.count {
		r.Clear()
		r.StartTime = r.time()
		r.StartBucket = 0
	}
	return r.Buckets[r.bucket(diff)]
}

func (r *rolling) bucket(i int) int {
	return (r.StartBucket + i) % r.count
}

func (r *rolling) Clear() {
	for i := range r.Buckets {
		r.Buckets[i].clear()
	}
}

func (r *rolling) Succeed() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.current()
	b.attempts += 1
}

func (r *rolling) Fail() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.current()
	b.attempts += 1
	b.failures += 1
}

// TODO: Calculate failure percent over window
// TODO: test success / failure recording over time
