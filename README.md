# Cacher

Inspired by the [Laravel Cache Facade](https://laravel.com/docs/10.x/cache), `cacher` provides a convenient wrapper around the official [Redis client](https://github.com/redis/go-redis). This library provides two methods for interacting with Redis. The first is via standard types such as strings, bytes, numbers, etc. However, often times we want to store structs. Cacher also ships with an Entity cache that uses generics. This client automatically marshalls the structs using JSON to and from the redis cache.

## Usage

## Client

The client provides a standard wrapper around the Redis client.

```golang

rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})

cache := cacher.New(rdb)

// fetch a value from the cache or from our database if it doesnt existing the cache
value, err := cache.RememberString(ctx, "my-key", time.Hour*24, func(ctx context.Context) (string, error) {
    // ... fetch value from database

    return databaseValue, nil
})

// put a value in the cache
err := cache.Put(ctx, "my-key", "hello-world", time.Hour*24)

err := cache.PutForever(ctx, "my-key", "hello-world")

// get a string from the cache
value, err := cache.GetString(ctx, "my-key")

// delete a key in the database
err := cache.Forget(ctx, "my-key")

```

## Entity Client

The entity client uses generics and JSON marshalling for automatically marhsalling data to and from the cache.

```golang
type MyEntity struct {}

rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})

cache := cacher.NewEntity[MyEntity](rdb)

// fetch a value from the cache or from our database if it doesnt existing the cache
value, err := cache.Remember(ctx, "my-key", time.Hour*24, func(ctx context.Context) (*MyEntity, error) {
    // ... fetch value from database
    return databaseValue, nil
})

// put a value in the cache
err := cache.Put(ctx, "my-key", &MyEntity{}, time.Hour*24)

err := cache.PutForever(ctx, "my-key", &MyEntity{})

// get a string from the cache
value, err := cache.Get(ctx, "my-key")

// delete a key in the database
err := cache.Forget(ctx, "my-key")
```
