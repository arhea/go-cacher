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

		value, err = client.RememberInt(ctx, key, time.Minute*5, func(ctx context.Context) (int, error) {
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

		value, err = client.RememberInt64(ctx, key, time.Minute*5, func(ctx context.Context) (int64, error) {
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

		value, err = client.RememberFloat32(ctx, key, time.Minute*5, func(ctx context.Context) (float32, error) {
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

		value, err = client.RememberFloat64(ctx, key, time.Minute*5, func(ctx context.Context) (float64, error) {
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

	t.Run("GetWithNotFoundError", func(t *testing.T) {
		_, err := client.GetString(ctx, "some-key")
		assert.ErrorIs(t, err, cacher.NotFoundError)
	})
}
