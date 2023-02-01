package rlock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/xhigher/hzgo/logger"
	"io"
	"sync/atomic"
	"time"
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("rlock: not obtained")
)

type SimpleLock struct {
	client *redis.Client
	ctx    context.Context
	key    string
	ttl    time.Duration
}

func NewSimpleLock(client *redis.Client, key string, ttl time.Duration) (lock *SimpleLock, err error) {
	ok, err := client.SetNX(context.Background(), key, "1", ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotObtained
	}
	return &SimpleLock{
		client: client,
		key:    key,
		ttl:    ttl,
	}, nil
}

func (c *SimpleLock) Release() (err error) {
	err = c.client.Del(context.Background(), c.key).Err()
	if err != nil {
		logger.Errorf("error %v", err)
	}
	return
}

// --------------------------------------------------------------------

// Lock represents an obtained, distributed lock.
type RetryLock struct {
	client *redis.Client
	key    string
	value  string
}

func NewRetryLock(client *redis.Client, key string, ttl time.Duration, retryDuration time.Duration, limit int) (*RetryLock, error) {
	ctx := context.Background()
	// Create a random token
	tmp := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, tmp); err != nil {
		return nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(tmp)

	option := &Options{
		RetryStrategy: LimitRetry(LinearBackoff(retryDuration), limit),
		Metadata:      "",
	}

	value := token + option.getMetadata()
	retry := option.getRetryStrategy()

	// make sure we don't retry forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(ttl))
		defer cancel()
	}

	var timer *time.Timer
	for {
		ok, err := client.SetNX(ctx, key, value, ttl).Result()
		if err != nil {
			return nil, err
		} else if ok {
			return &RetryLock{client: client, key: key, value: value}, nil
		}

		backoff := retry.NextBackoff()
		if backoff < 1 {
			return nil, ErrNotObtained
		}

		if timer == nil {
			timer = time.NewTimer(backoff)
			defer timer.Stop()
		} else {
			timer.Reset(backoff)
		}

		select {
		case <-ctx.Done():
			return nil, ErrNotObtained
		case <-timer.C:
		}
	}
}

func (l *RetryLock) IsOwn() (own bool, err error) {
	re, err := l.client.Get(context.Background(), l.key).Result()
	if err == redis.Nil {
		return false, ErrNotObtained
	} else if err != nil {
		return false, err
	}
	if re != l.value {
		return false, nil
	}
	return true, nil
}

// Key returns the redis key used by the lock.
func (l *RetryLock) Key() string {
	return l.key
}

// Token returns the token value set by the lock.
func (l *RetryLock) Token() string {
	return l.value[:22]
}

// Metadata returns the metadata of the lock.
func (l *RetryLock) Metadata() string {
	return l.value[22:]
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (l *RetryLock) TTL() (time.Duration, error) {
	return l.client.PTTL(context.Background(), l.key).Result()
}

// Refresh extends the lock with a new TTL.
// May return ErrNotObtained if refresh is unsuccessful.
func (l *RetryLock) Refresh(ttl time.Duration) error {
	return l.client.Expire(context.Background(), l.key, ttl).Err()
}

// Release manually releases the lock.
// May return ErrLockNotHeld.
func (l *RetryLock) Release() error {
	return l.client.Del(context.Background(), l.key).Err()
}

// --------------------------------------------------------------------

// Options describe the options for the lock
type Options struct {
	// RetryStrategy allows to customise the lock retry strategy.
	// Default: do not retry
	RetryStrategy RetryStrategy

	// Metadata string is appended to the lock token.
	Metadata string
}

func (o *Options) getMetadata() string {
	if o != nil {
		return o.Metadata
	}
	return ""
}

func (o *Options) getRetryStrategy() RetryStrategy {
	if o != nil && o.RetryStrategy != nil {
		return o.RetryStrategy
	}
	return NoRetry()
}

// --------------------------------------------------------------------

// RetryStrategy allows to customise the lock retry strategy.
type RetryStrategy interface {
	// NextBackoff returns the next backoff duration.
	NextBackoff() time.Duration
}

type linearBackoff time.Duration

// LinearBackoff allows retries regularly with customized intervals
func LinearBackoff(backoff time.Duration) RetryStrategy {
	return linearBackoff(backoff)
}

// NoRetry acquire the lock only once.
func NoRetry() RetryStrategy {
	return linearBackoff(0)
}

func (r linearBackoff) NextBackoff() time.Duration {
	return time.Duration(r)
}

type limitedRetry struct {
	s   RetryStrategy
	cnt int64
	max int64
}

// LimitRetry limits the number of retries to max attempts.
func LimitRetry(s RetryStrategy, max int) RetryStrategy {
	return &limitedRetry{s: s, max: int64(max)}
}

func (r *limitedRetry) NextBackoff() time.Duration {
	if atomic.LoadInt64(&r.cnt) >= r.max {
		return 0
	}
	atomic.AddInt64(&r.cnt, 1)
	return r.s.NextBackoff()
}

type exponentialBackoff struct {
	cnt uint64

	min, max time.Duration
}

// ExponentialBackoff strategy is an optimization strategy with a retry time of 2**n milliseconds (n means number of times).
// You can set a minimum and maximum value, the recommended minimum value is not less than 16ms.
func ExponentialBackoff(min, max time.Duration) RetryStrategy {
	return &exponentialBackoff{min: min, max: max}
}

func (r *exponentialBackoff) NextBackoff() time.Duration {
	cnt := atomic.AddUint64(&r.cnt, 1)

	ms := 2 << 25
	if cnt < 25 {
		ms = 2 << cnt
	}

	if d := time.Duration(ms) * time.Millisecond; d < r.min {
		return r.min
	} else if r.max != 0 && d > r.max {
		return r.max
	} else {
		return d
	}
}


