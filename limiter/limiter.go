package limiter

import (
	"github.com/hashicorp/golang-lru"
	"github.com/juju/ratelimit"
	"sync"
	"time"
)

type Limiter struct {
	config          Config
	generatorBucket *ratelimit.Bucket
	tokenCache      *lru.Cache
	tokenLock       sync.Mutex
	hostCache       *lru.Cache
	hostLock        sync.Mutex
}

type Config struct {
	GeneratorBucketCapacity int64
	GeneratorBucketInterval time.Duration

	// the token limiter cache capacity
	TokenCacheSize      int
	TokenBucketCapacity int64
	TokenBucketInterval time.Duration

	// the host limiter cache capacity
	HostCacheSize      int
	HostBucketCapacity int64
	HostBucketInterval time.Duration
}

func New(c Config) (*Limiter, error) {
	tokenCache, err := lru.New(c.TokenCacheSize)
	if err != nil {
		return nil, err
	}

	hostCache, err := lru.New(c.HostCacheSize)
	if err != nil {
		return nil, err
	}

	return &Limiter{
		config: c,

		// the token bucket rate-limits global token generation.
		// it refills at a rate of 1 / second,
		// with an initial capacity of 120 tokens
		generatorBucket: ratelimit.NewBucket(c.GeneratorBucketInterval, c.GeneratorBucketCapacity),

		// holds buckets for rate-limiting individual tokens
		tokenCache: tokenCache,

		// holds buckets for rate-limiting individual hosts
		hostCache: hostCache,
	}, nil
}

// LimitGenerator is a rate-limiter for global token generation.
// It immediately returns a duration until tokens are available.
func (l *Limiter) LimitGenerator() int {
	return int(l.generatorBucket.Take(1).Seconds())
}

// LimitToken is a rate-limiter for how often individual tokens
// can be used. It immediately returns a duration until tokens are available.
func (l *Limiter) LimitToken(key string) int {
	var buck *ratelimit.Bucket

	l.tokenLock.Lock()
	defer l.tokenLock.Unlock()

	val, ok := l.tokenCache.Get(key)
	if ok {
		buck = val.(*ratelimit.Bucket)
	} else {
		buck = ratelimit.NewBucket(l.config.TokenBucketInterval, l.config.TokenBucketCapacity)
	}

	l.tokenCache.Add(key, buck)
	return int(buck.Take(1).Seconds())
}

// LimitHost is a rate-limiter for how often individual hosts
// can be used. It immediately returns a duration until a host is available.
func (l *Limiter) LimitHost(key string) int {
	var buck *ratelimit.Bucket

	l.hostLock.Lock()
	defer l.hostLock.Unlock()

	val, ok := l.hostCache.Get(key)
	if ok {
		buck = val.(*ratelimit.Bucket)
	} else {
		buck = ratelimit.NewBucket(l.config.HostBucketInterval, l.config.HostBucketCapacity)
	}

	l.hostCache.Add(key, buck)
	return int(buck.Take(1).Seconds())
}
