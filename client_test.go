package cacher_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/arhea/go-cacher"
	mockredis "github.com/arhea/go-mock-redis"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock, err := mockredis.NewClient(ctx, t)

	if err != nil {
		t.Fatal(err)
		return
	}

	r := mock.Client()
	client := cacher.New(r)

	t.Run("Basic", func(t *testing.T) {
		key := "basic"
		value := "hello-world"

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.Get(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value, fetchedValue, "should be equal as the value was stored")

		err = client.Forget(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err = client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue, "should be false as the key has been forgotten")
	})

	t.Run("BasicString", func(t *testing.T) {
		key := "basic-string"
		value := "hello-world"

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.GetString(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value, fetchedValue, "should be equal as the value was stored")

		err = client.Forget(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err = client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue, "should be false as the key has been forgotten")
	})

	t.Run("GetStringWithDefault", func(t *testing.T) {
		key := "default-string"

		cacheValue := "cache-value"
		defaultValue := "default-value"

		value1 := client.GetStringWithDefault(ctx, key, defaultValue)

		assert.Equal(t, defaultValue, value1, "should be equal as the value was not stored")

		err := client.Put(ctx, key, cacheValue, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		value2 := client.GetStringWithDefault(ctx, key, defaultValue)

		assert.Equal(t, cacheValue, value2, "should be equal as the value is stored")
	})

	t.Run("GetIntWithDefault", func(t *testing.T) {
		key := "default-int"

		cacheValue := int(92)
		defaultValue := int(23)

		value1 := client.GetIntWithDefault(ctx, key, defaultValue)

		assert.Equal(t, defaultValue, value1, "should be equal as the value was not stored")

		err := client.Put(ctx, key, cacheValue, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		value2 := client.GetIntWithDefault(ctx, key, defaultValue)

		assert.Equal(t, cacheValue, value2, "should be equal as the value is stored")
	})

	t.Run("GetInt64WithDefault", func(t *testing.T) {
		key := "default-int64"

		cacheValue := int64(92)
		defaultValue := int64(23)

		value1 := client.GetInt64WithDefault(ctx, key, defaultValue)

		assert.Equal(t, defaultValue, value1, "should be equal as the value was not stored")

		err := client.Put(ctx, key, cacheValue, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		value2 := client.GetInt64WithDefault(ctx, key, defaultValue)

		assert.Equal(t, cacheValue, value2, "should be equal as the value is stored")
	})

	t.Run("GetFloat32WithDefault", func(t *testing.T) {
		key := "default-float32"

		cacheValue := float32(92)
		defaultValue := float32(23)

		value1 := client.GetFloat32WithDefault(ctx, key, defaultValue)

		assert.Equal(t, defaultValue, value1, "should be equal as the value was not stored")

		err := client.Put(ctx, key, cacheValue, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		value2 := client.GetFloat32WithDefault(ctx, key, defaultValue)

		assert.Equal(t, cacheValue, value2, "should be equal as the value is stored")
	})

	t.Run("GetFloat64WithDefault", func(t *testing.T) {
		key := "default-float64"

		cacheValue := float64(92)
		defaultValue := float64(23)

		value1 := client.GetFloat64WithDefault(ctx, key, defaultValue)

		assert.Equal(t, defaultValue, value1, "should be equal as the value was not stored")

		err := client.Put(ctx, key, cacheValue, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		value2 := client.GetFloat64WithDefault(ctx, key, defaultValue)

		assert.Equal(t, cacheValue, value2, "should be equal as the value is stored")
	})

	t.Run("GetBoolWithDefault", func(t *testing.T) {
		key := "default-bool"

		cacheValue := true
		defaultValue := false

		value1 := client.GetBoolWithDefault(ctx, key, defaultValue)

		assert.Equal(t, defaultValue, value1, "should be equal as the value was not stored")

		err := client.Put(ctx, key, cacheValue, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		value2 := client.GetBoolWithDefault(ctx, key, defaultValue)

		assert.Equal(t, cacheValue, value2, "should be equal as the value is stored")
	})

	t.Run("GetBytesWithDefault", func(t *testing.T) {
		key := "default-bytes"

		cacheValue := []byte("hello-world")
		defaultValue := []byte("hello-mars")

		value1 := client.GetBytesWithDefault(ctx, key, defaultValue)

		assert.Equal(t, defaultValue, value1, "should be equal as the value was not stored")

		err := client.Put(ctx, key, cacheValue, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		value2 := client.GetBytesWithDefault(ctx, key, defaultValue)

		assert.Equal(t, cacheValue, value2, "should be equal as the value is stored")
	})

	t.Run("RememberString", func(t *testing.T) {
		key := "remember-string"
		value1 := "hello-world-1"

		value, err := client.RememberString(ctx, key, time.Minute*5, func(ctx context.Context) (string, error) {
			return value1, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should be equal to the fetched value")

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		value, err = client.RememberString(ctx, key, time.Minute*5, func(ctx context.Context) (string, error) {
			return "some other value", nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should still be equal to the fetched value as it should be in the cache")
	})

	t.Run("BasicBool", func(t *testing.T) {
		key := "basic-bool"
		value := true

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.GetBool(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value, fetchedValue, "should be equal as the value was stored")
	})

	t.Run("RememberBool", func(t *testing.T) {
		key := "remember-bool"
		value1 := true

		value, err := client.RememberBool(ctx, key, time.Minute*5, func(ctx context.Context) (bool, error) {
			return value1, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should be equal to the fetched value")

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		value, err = client.RememberBoolForever(ctx, key, func(ctx context.Context) (bool, error) {
			return false, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should still be equal to the fetched value as it should be in the cache")
	})

	t.Run("BasicInt", func(t *testing.T) {
		key := "basic-int"
		value := int(12)

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.GetInt(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value, fetchedValue, "should be equal as the value was stored")
	})

	t.Run("RememberInt", func(t *testing.T) {
		key := "remember-int"
		value1 := int(7)

		value, err := client.RememberInt(ctx, key, time.Minute*5, func(ctx context.Context) (int, error) {
			return value1, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should be equal to the fetched value")

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		value, err = client.RememberIntForever(ctx, key, func(ctx context.Context) (int, error) {
			return 47, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should still be equal to the fetched value as it should be in the cache")
	})

	t.Run("BasicInt64", func(t *testing.T) {
		key := "basic-int64"
		value := int64(12)

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.GetInt64(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value, fetchedValue, "should be equal as the value was stored")
	})

	t.Run("RememberInt64", func(t *testing.T) {
		key := "remember-int64"
		value1 := int64(7)

		value, err := client.RememberInt64(ctx, key, time.Minute*5, func(ctx context.Context) (int64, error) {
			return value1, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should be equal to the fetched value")

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		value, err = client.RememberInt64Forever(ctx, key, func(ctx context.Context) (int64, error) {
			return 47, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should still be equal to the fetched value as it should be in the cache")
	})

	t.Run("BasicIntFloat32", func(t *testing.T) {
		key := "basic-float32"
		value := float32(12.1)

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.GetFloat32(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value, fetchedValue, "should be equal as the value was stored")
	})

	t.Run("RememberFloat32", func(t *testing.T) {
		key := "remember-float32"
		value1 := float32(7.7)

		value, err := client.RememberFloat32(ctx, key, time.Minute*5, func(ctx context.Context) (float32, error) {
			return value1, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should be equal to the fetched value")

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		value, err = client.RememberFloat32Forever(ctx, key, func(ctx context.Context) (float32, error) {
			return 47, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should still be equal to the fetched value as it should be in the cache")
	})

	t.Run("BasicIntFloat64", func(t *testing.T) {
		key := "basic-float64"
		value := float64(12.1)

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.GetFloat64(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value, fetchedValue, "should be equal as the value was stored")
	})

	t.Run("RememberFloat64", func(t *testing.T) {
		key := "remember-float64"
		value1 := float64(7.7)

		value, err := client.RememberFloat64(ctx, key, time.Minute*5, func(ctx context.Context) (float64, error) {
			return value1, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should be equal to the fetched value")

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		value, err = client.RememberFloat64Forever(ctx, key, func(ctx context.Context) (float64, error) {
			return 47, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should still be equal to the fetched value as it should be in the cache")
	})

	t.Run("BasicBytes", func(t *testing.T) {
		key := "basic-float64"
		value := []byte("hello-world")

		err := client.Put(ctx, key, value, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		fetchedValue, err := client.GetBytes(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, bytes.Equal(value, fetchedValue), "should be equal as the value was stored")
	})

	t.Run("RememberBytes", func(t *testing.T) {
		key := "remember-bytes"
		value1 := []byte("hello-world")

		value, err := client.RememberBytes(ctx, key, time.Minute*5, func(ctx context.Context) ([]byte, error) {
			return value1, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should be equal to the fetched value")

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, hasValue, "should be true as the key exists")

		value, err = client.RememberBytesForever(ctx, key, func(ctx context.Context) ([]byte, error) {
			return []byte("hello-cruel-world"), nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, value1, value, "value should still be equal to the fetched value as it should be in the cache")
	})

	t.Run("TestTTL", func(t *testing.T) {
		key := "test-ttl"
		value := "hello-world"

		err := client.Put(ctx, key, value, time.Second*1)

		if err != nil {
			t.Error(err)
			return
		}

		time.Sleep(time.Second * 2)

		hasValue, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue, "should be expired and therefore not exist")

	})

	t.Run("ForgetMany", func(t *testing.T) {
		_ = client.Put(ctx, "prefix-1", "something", time.Minute*5)
		_ = client.Put(ctx, "prefix-2", "something", time.Minute*5)
		_ = client.Put(ctx, "prefix-3", "something", time.Minute*5)
		_ = client.Put(ctx, "prefix-4", "something", time.Minute*5)

		err := client.ForgetWithPrefix(ctx, "prefix-*")

		if err != nil {
			t.Error(err)
			return
		}

		hasValue1, err := client.Has(ctx, "prefix-1")

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue1, "should be false as the key has been forgotten")

		hasValue2, err := client.Has(ctx, "prefix-2")

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue2, "should be false as the key has been forgotten")

		hasValue3, err := client.Has(ctx, "prefix-3")

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue3, "should be false as the key has been forgotten")

		hasValue4, err := client.Has(ctx, "prefix-4")

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue4, "should be false as the key has been forgotten")

		hasValue5, err := client.Has(ctx, "prefix-5")

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, hasValue5, "should be false as the key has been forgotten")
	})

	t.Run("IncrementDecrement", func(t *testing.T) {
		key := "counter"

		err := client.Put(ctx, key, 1, time.Minute*5)

		if err != nil {
			t.Error(err)
			return
		}

		err = client.Increment(ctx, key, 5)

		if err != nil {
			t.Error(err)
			return
		}

		value1, err := client.GetInt(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 6, value1, "should be equal as the value was incremented")

		err = client.Decrement(ctx, key, 1)

		if err != nil {
			t.Error(err)
			return
		}

		value2, err := client.GetInt(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 5, value2, "should be equal as the value was decremented")

	})

	t.Run("GetWithNotFoundError", func(t *testing.T) {
		_, err := client.GetString(ctx, "some-key")
		assert.ErrorIs(t, err, cacher.NotFoundError)
	})
}
