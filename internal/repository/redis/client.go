package redis

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client  *redis.Client
	Redsync *redsync.Redsync
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return &RedisClient{
		Client:  client,
		Redsync: rs,
	}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// DistributedLock provides distributed locking using Redsync
type DistributedLock struct {
	mu   *redsync.Mutex
	name string
	rs   *redsync.Redsync
}

func (r *RedisClient) NewLock(name string, options ...redsync.Option) *DistributedLock {
	return &DistributedLock{
		mu:   r.Redsync.NewMutex(name, options...),
		name: name,
		rs:   r.Redsync,
	}
}

func (l *DistributedLock) Lock(ctx context.Context) error {
	return l.mu.LockContext(ctx)
}

func (l *DistributedLock) Unlock(ctx context.Context) (bool, error) {
	return l.mu.UnlockContext(ctx)
}

func (l *DistributedLock) Extend(ctx context.Context) (bool, error) {
	return l.mu.ExtendContext(ctx)
}

// Lock options for common use cases
func WithExpiry(expiry time.Duration) redsync.Option {
	return redsync.WithExpiry(expiry)
}

func WithTries(tries int) redsync.Option {
	return redsync.WithTries(tries)
}

func WithRetryDelay(delay time.Duration) redsync.Option {
	return redsync.WithRetryDelay(delay)
}

func WithDriftFactor(factor float64) redsync.Option {
	return redsync.WithDriftFactor(factor)
}
