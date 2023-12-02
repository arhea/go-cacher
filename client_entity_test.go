package cacher_test

import (
	"context"
	"testing"
	"time"

	cacher "github.com/arhea/go-cacher"
	mockredis "github.com/arhea/go-mock-redis"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

type TestEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestEntityClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock, err := mockredis.NewClient(ctx, t)

	if err != nil {
		t.Fatal(err)
		return
	}

	r := mock.Client()
	client := cacher.NewEntity[TestEntity](r)

	t.Run("BasicFunctions", func(t *testing.T) {
		entity := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		err := client.Put(ctx, entity.ID, entity, time.Minute*1)

		if err != nil {
			t.Error(err)
			return
		}

		has, err := client.Has(ctx, entity.ID)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, has)

		fetchedEntity, err := client.Get(ctx, entity.ID)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, fetchedEntity.ID, entity.ID)
		assert.Equal(t, fetchedEntity.Name, entity.Name)

		err = client.Forget(ctx, entity.ID)

		if err != nil {
			t.Error(err)
			return
		}

		has, err = client.Has(ctx, entity.ID)

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, has)
	})
}
