package cacher

import "errors"

// NotFoundError is returned when a key is not found in the cache.
var NotFoundError = errors.New("key not found")

// EntityMarshalError is returned when an entity cannot be marshalled or unmarshalled.
var EntityMarshalError = errors.New("error marshalling entity")
