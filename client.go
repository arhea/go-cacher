package cacher

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// New creates a new instance of the Cache client from an existing redis client. This will not close the
// redis client.
func New(r *redis.Client) *Client {
	return &Client{
		redis: r,
	}
}

// Client is a client that simplifies the access to the redis for common caching patterns.
type Client struct {
	redis *redis.Client
}

// Has checks if a key exists in the cache. It returns false if the key does not exist or there was an error. It will
// also return the error if there is one.
func (c *Client) Has(ctx context.Context, key string) (bool, error) {
	cmd := c.redis.Exists(ctx, key)
	err := cmd.Err()

	if err != nil {
		return false, err
	}

	return cmd.Val() == 1, nil
}

// Forget removes a key from the cache. It returns an error if there was one.
func (c *Client) Forget(ctx context.Context, key string) error {
	cmd := c.redis.Del(ctx, key)
	return cmd.Err()
}

// ForgetWithPrefix removes all keys from the cache that match the given prefix. It returns an error if there was one.
func (c *Client) ForgetWithPrefix(ctx context.Context, prefix string) error {
	iter := c.redis.Scan(ctx, 0, prefix, 0).Iterator()

	for iter.Next(ctx) {
		if err := c.Forget(ctx, iter.Val()); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

// Put adds a value to the cache with an expiration. It returns an error if there was one.
func (c *Client) Put(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	cmd := c.redis.Set(ctx, key, value, exp)
	return cmd.Err()
}

// PutForever adds a value to the cache without an expiration. It returns an error if there was one.
func (c *Client) PutForever(ctx context.Context, key string, value interface{}) error {
	return c.Put(ctx, key, value, 0)
}

// Increment increments a value in the cache. It returns an error if there was one.
func (c *Client) Increment(ctx context.Context, key string, value int64) error {
	cmd := c.redis.IncrBy(ctx, key, value)
	return cmd.Err()
}

// Decrement decrements a value in the cache. It returns an error if there was one.
func (c *Client) Decrement(ctx context.Context, key string, value int64) error {
	cmd := c.redis.DecrBy(ctx, key, value)
	return cmd.Err()
}

// Get retrieves a value from the cache. It returns an error if there was one. If the key does not exist it will return
// a NotFoundError.
func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return cmd.Val(), nil
}

// GetString returns the key as a string. If there was an error it will return an error. If the key does not exist it
// will return a NotFoundError. If there was an error the string value will be a zero string.
func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return "", NotFoundError
	}

	if err != nil {
		return "", err
	}

	return cmd.Val(), nil
}

// GetBytes returns the key as a []byte. If there was an error it will return an error. If the key does not exist it
// will return a NotFoundError. If there was an error the value will be nil.
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return nil, NotFoundError
	}

	if err != nil {
		return nil, err
	}

	return cmd.Bytes()
}

// GetBool returns the key as a bool. If there was an error it will return an error. If the key does not exist it
// will return a NotFoundError. If there was an error the value will be false.
func (c *Client) GetBool(ctx context.Context, key string) (bool, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return false, NotFoundError
	}

	if err != nil {
		return false, err
	}

	return cmd.Bool()
}

// GetInt returns the key as an int. If there was an error it will return an error. If the key does not exist it
// will return a NotFoundError. If there was an error the value will be 0.
func (c *Client) GetInt(ctx context.Context, key string) (int, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return 0, NotFoundError
	}

	if err != nil {
		return 0, err
	}

	return cmd.Int()
}

// GetInt64 returns the key as an int64. If there was an error it will return an error. If the key does not exist it
// will return a NotFoundError. If there was an error the value will be 0.
func (c *Client) GetInt64(ctx context.Context, key string) (int64, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return 0, NotFoundError
	}

	if err != nil {
		return 0, err
	}

	return cmd.Int64()
}

// GetFloat32 returns the key as an float32. If there was an error it will return an error. If the key does not exist it
// will return a NotFoundError. If there was an error the value will be 0.
func (c *Client) GetFloat32(ctx context.Context, key string) (float32, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return 0, NotFoundError
	}

	if err != nil {
		return 0, err
	}

	return cmd.Float32()
}

// GetFloat64 returns the key as an float64. If there was an error it will return an error. If the key does not exist it
// will return a NotFoundError. If there was an error the value will be 0.
func (c *Client) GetFloat64(ctx context.Context, key string) (float64, error) {
	cmd := c.redis.Get(ctx, key)
	err := cmd.Err()

	if errors.Is(err, redis.Nil) {
		return 0, NotFoundError
	}

	if err != nil {
		return 0, err
	}

	return cmd.Float64()
}

// GetStringWithDefault will return the value as a string. If there was an error or the value is zero, it will return
// the default value.
func (c *Client) GetStringWithDefault(ctx context.Context, key string, defaultValue string) string {
	val, err := c.GetString(ctx, key)

	if err != nil || val == "" {
		return defaultValue
	}

	return val
}

// GetBoolWithDefault will return the value as a string. If there was an error or the value is zero, it will return
// the default value.
func (c *Client) GetBoolWithDefault(ctx context.Context, key string, defaultValue bool) bool {
	val, err := c.GetBool(ctx, key)

	if err != nil {
		return defaultValue
	}

	return val
}

// GetBytesWithDefault will return the value as a []byte. If there was an error or the value is nil, it will return
// the default value.
func (c *Client) GetBytesWithDefault(ctx context.Context, key string, defaultValue []byte) []byte {
	val, err := c.GetBytes(ctx, key)

	if err != nil || val == nil {
		return defaultValue
	}

	return val
}

// GetIntWithDefault will return the value as a int. If there was an error or the value is zero, it will return
// the default value.
func (c *Client) GetIntWithDefault(ctx context.Context, key string, defaultValue int) int {
	val, err := c.GetInt(ctx, key)

	if err != nil || val == 0 {
		return defaultValue
	}

	return val
}

// GetInt64WithDefault will return the value as a int64. If there was an error or the value is zero, it will return
// the default value.
func (c *Client) GetInt64WithDefault(ctx context.Context, key string, defaultValue int64) int64 {
	val, err := c.GetInt64(ctx, key)

	if err != nil || val == 0 {
		return defaultValue
	}

	return val
}

// GetFloat32WithDefault will return the value as a float32. If there was an error or the value is zero, it will return
// the default value.
func (c *Client) GetFloat32WithDefault(ctx context.Context, key string, defaultValue float32) float32 {
	val, err := c.GetFloat32(ctx, key)

	if err != nil || val == 0 {
		return defaultValue
	}

	return val
}

// GetFloat64WithDefault will return the value as a float64. If there was an error or the value is zero, it will return
// the default value.
func (c *Client) GetFloat64WithDefault(ctx context.Context, key string, defaultValue float64) float64 {
	val, err := c.GetFloat64(ctx, key)

	if err != nil || val == 0 {
		return defaultValue
	}

	return val
}

// RememberString will attempt to get the value from the cache. If it is not found it will call the fetcher function to
// get the value. It will then put the value in the cache for later. It returns the value or an error if there was one.
func (c *Client) RememberString(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) (string, error)) (string, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetString(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && val != "" {
		return val, nil
	}

	// if there was an error and it wasn't a not found error return
	if err != nil && !errors.Is(err, NotFoundError) {
		return "", err
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return "", err
	}

	// put the value in the cache for later
	if err := c.Put(ctx, key, val, exp); err != nil {
		return "", err
	}

	return val, nil
}

// RememberStringForever is the same as RememberString but it will not expire the value.
func (c *Client) RememberStringForever(ctx context.Context, key string, fetcher func(ctx context.Context) (string, error)) (string, error) {
	return c.RememberString(ctx, key, 0, fetcher)
}

// RememberBool will attempt to get the value from the cache. If it is not found it will call the fetcher function to
// get the value. It will then put the value in the cache for later. It returns the value or an error if there was one.
func (c *Client) RememberBool(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) (bool, error)) (bool, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetBool(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil {
		return val, nil
	}

	// if there was an error and it wasn't a not found error return
	if err != nil && !errors.Is(err, NotFoundError) {
		return false, err
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return false, err
	}

	// put the value in the cache for later
	if err := c.Put(ctx, key, val, exp); err != nil {
		return false, err
	}

	return val, nil
}

// RememberBoolForever is the same as RememberBool but it will not expire the value.
func (c *Client) RememberBoolForever(ctx context.Context, key string, fetcher func(ctx context.Context) (bool, error)) (bool, error) {
	return c.RememberBool(ctx, key, 0, fetcher)
}

// RememberBytes will attempt to get the value from the cache. If it is not found it will call the fetcher function to
// get the value. It will then put the value in the cache for later. It returns the value or an error if there was one.
func (c *Client) RememberBytes(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) ([]byte, error)) ([]byte, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetBytes(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && len(val) > 0 {
		return val, nil
	}

	// if there was an error and it wasn't a not found error return
	if err != nil && !errors.Is(err, NotFoundError) {
		return nil, err
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return nil, err
	}

	// put the value in the cache for later
	if err := c.Put(ctx, key, val, exp); err != nil {
		return nil, err
	}

	return val, nil
}

// RememberBytesForever is the same as RememberString but it will not expire the value.
func (c *Client) RememberBytesForever(ctx context.Context, key string, fetcher func(ctx context.Context) ([]byte, error)) ([]byte, error) {
	return c.RememberBytes(ctx, key, 0, fetcher)
}

// RememberInt will attempt to get the value from the cache. If it is not found it will call the fetcher function to
// get the value. It will then put the value in the cache for later. It returns the value or an error if there was one.
func (c *Client) RememberInt(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) (int, error)) (int, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetInt(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && val != 0 {
		return val, nil
	}

	// if there was an error and it wasn't a not found error return
	if err != nil && !errors.Is(err, NotFoundError) {
		return 0, err
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return 0, err
	}

	// put the value in the cache for later
	if err := c.Put(ctx, key, val, exp); err != nil {
		return 0, err
	}

	return val, nil
}

// RememberIntForever is the same as RememberInt but it will not expire the value.
func (c *Client) RememberIntForever(ctx context.Context, key string, fetcher func(ctx context.Context) (int, error)) (int, error) {
	return c.RememberInt(ctx, key, 0, fetcher)
}

// RememberInt64 will attempt to get the value from the cache. If it is not found it will call the fetcher function to
// get the value. It will then put the value in the cache for later. It returns the value or an error if there was one.
func (c *Client) RememberInt64(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) (int64, error)) (int64, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetInt64(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && val != 0 {
		return val, nil
	}

	// if there was an error and it wasn't a not found error return
	if err != nil && !errors.Is(err, NotFoundError) {
		return 0, err
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return 0, err
	}

	// put the value in the cache for later
	if err := c.Put(ctx, key, val, exp); err != nil {
		return 0, err
	}

	return val, nil
}

// RememberInt64Forever is the same as RememberInt64 but it will not expire the value.
func (c *Client) RememberInt64Forever(ctx context.Context, key string, fetcher func(ctx context.Context) (int64, error)) (int64, error) {
	return c.RememberInt64(ctx, key, 0, fetcher)
}

// RememberFloat32 will attempt to get the value from the cache. If it is not found it will call the fetcher function to
// get the value. It will then put the value in the cache for later. It returns the value or an error if there was one.
func (c *Client) RememberFloat32(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) (float32, error)) (float32, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetFloat32(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && val != 0 {
		return val, nil
	}

	// if there was an error and it wasn't a not found error return
	if err != nil && !errors.Is(err, NotFoundError) {
		return 0, err
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return 0, err
	}

	// put the value in the cache for later
	if err := c.Put(ctx, key, val, exp); err != nil {
		return 0, err
	}

	return val, nil
}

// RememberFloat32Forever is the same as RememberFloat32 but it will not expire the value.
func (c *Client) RememberFloat32Forever(ctx context.Context, key string, fetcher func(ctx context.Context) (float32, error)) (float32, error) {
	return c.RememberFloat32(ctx, key, 0, fetcher)
}

// RememberFloat64 will attempt to get the value from the cache. If it is not found it will call the fetcher function to
// get the value. It will then put the value in the cache for later. It returns the value or an error if there was one.
func (c *Client) RememberFloat64(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) (float64, error)) (float64, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetFloat64(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && val != 0 {
		return val, nil
	}

	// if there was an error and it wasn't a not found error return
	if err != nil && !errors.Is(err, NotFoundError) {
		return 0, err
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return 0, err
	}

	// put the value in the cache for later
	if err := c.Put(ctx, key, val, exp); err != nil {
		return 0, err
	}

	return val, nil
}

// RememberFloat64Forever is the same as RememberFloat64 but it will not expire the value.
func (c *Client) RememberFloat64Forever(ctx context.Context, key string, fetcher func(ctx context.Context) (float64, error)) (float64, error) {
	return c.RememberFloat64(ctx, key, 0, fetcher)
}
