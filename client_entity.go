package cacher

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewEntityWithClient creates a new EntityClient with the given Client.
func NewEntityWithClient[E any](c *Client) *EntityClient[E] {
	return &EntityClient[E]{
		client: c,
	}
}

// NewEntity creates a new EntityClient with the given redis client.
func NewEntity[E any](r *redis.Client) *EntityClient[E] {
	return NewEntityWithClient[E](New(r))
}

// EntityClient is a wrapper around the Client that provides a more convenient access parttern using generics. This
// client will automatically marshal and unmarshal entities to and from the cache using JSON.
type EntityClient[E any] struct {
	client *Client
}

// Has checks if the given key exists in the cache.
func (c *EntityClient[E]) Has(ctx context.Context, key string) (bool, error) {
	return c.client.Has(ctx, key)
}

// Forget removes the given key from the cache.
func (c *EntityClient[E]) Forget(ctx context.Context, key string) error {
	return c.client.Forget(ctx, key)
}

// ForgetWithPrefix removes all keys from the cache that start with the given prefix.
func (c *EntityClient[E]) ForgetWithPrefix(ctx context.Context, prefix string) error {
	return c.client.ForgetWithPrefix(ctx, prefix)
}

// Get fetches the entity from the cache. The value will be nil if it is not found along with an error.
func (c *EntityClient[E]) Get(ctx context.Context, key string) (*E, error) {
	// get the entity from the cache
	data, err := c.client.GetBytes(ctx, key)

	if err != nil {
		return nil, err
	}

	// unmarshal the entity
	var entity E

	err = json.Unmarshal(data, &entity)

	if err != nil {
		return nil, errors.Join(EntityMarshalError, err)
	}

	return &entity, err
}

// GetMany fetches the entities from the cache. The value will be nil if it is not found along with an error.
func (c *EntityClient[E]) GetMany(ctx context.Context, key string) ([]*E, error) {
	// get the entity from the cache
	data, err := c.client.GetBytes(ctx, key)

	if err != nil {
		return nil, err
	}

	// unmarshal the entity
	entities := make([]*E, 0)

	err = json.Unmarshal(data, &entities)

	if err != nil {
		return nil, errors.Join(EntityMarshalError, err)
	}

	return entities, err
}

// Put stores the entity in the cache for the given duration. If the duration is 0 the entity will be stored forever.
func (c *EntityClient[E]) Put(ctx context.Context, key string, value *E, exp time.Duration) error {
	// marshal the entity
	data, err := json.Marshal(value)

	if err != nil {
		return errors.Join(EntityMarshalError, err)
	}

	// put the entity into the cache
	return c.client.Put(ctx, key, data, exp)
}

// PutForever stores the entity in the cache forever.
func (c *EntityClient[E]) PutForever(ctx context.Context, key string, value *E) error {
	return c.Put(ctx, key, value, 0)
}

// PutMany stores the entities in the cache for the given duration. If the duration is 0 the entities will be stored forever.
func (c *EntityClient[E]) PutMany(ctx context.Context, key string, values []*E, exp time.Duration) error {
	// marshal the entity
	data, err := json.Marshal(values)

	if err != nil {
		return errors.Join(EntityMarshalError, err)
	}

	// put the entity into the cache
	return c.client.Put(ctx, key, data, exp)
}

// PutManyForever stores the entities in the cache forever.
func (c *EntityClient[E]) PutManyForever(ctx context.Context, key string, values []*E) error {
	return c.PutMany(ctx, key, values, 0)
}

// Remember fetches the entity from the cache if it exists. If it does not exist the fetcher will be called to get the
// entity and it will be stored in the cache for the given duration. If the duration is 0 the entity will be stored forever.
func (c *EntityClient[E]) Remember(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) (*E, error)) (*E, error) {
	// attempt to fetch the value from the cache
	val, err := c.Get(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && val != nil {
		return val, nil
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

// RememberForever wraps Remember and stores the entity in the cache forever.
func (c *EntityClient[E]) RememberForever(ctx context.Context, key string, fetcher func(ctx context.Context) (*E, error)) (*E, error) {
	return c.Remember(ctx, key, 0, fetcher)
}

// RememberMany fetches the entity from the cache if it exists. If it does not exist the fetcher will be called to get the
// entity and it will be stored in the cache for the given duration. If the duration is 0 the entity will be stored forever.
func (c *EntityClient[E]) RememberMany(ctx context.Context, key string, exp time.Duration, fetcher func(ctx context.Context) ([]*E, error)) ([]*E, error) {
	// attempt to fetch the value from the cache
	val, err := c.GetMany(ctx, key)

	// if there was no error and the value isn't blank return
	if err == nil && val != nil {
		return val, nil
	}

	// call the fetcher to get the value we should remember
	val, err = fetcher(ctx)

	if err != nil {
		return nil, err
	}

	// put the value in the cache for later
	if err := c.PutMany(ctx, key, val, exp); err != nil {
		return nil, err
	}

	return val, nil
}

// RememberManyForever wraps RememberMany and stores the entity in the cache forever.
func (c *EntityClient[E]) RememberManyForever(ctx context.Context, key string, fetcher func(ctx context.Context) ([]*E, error)) ([]*E, error) {
	return c.RememberMany(ctx, key, 0, fetcher)
}
