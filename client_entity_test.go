package cacher_test

import (
	"context"
	"fmt"
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

	t.Run("SingleFunctions", func(t *testing.T) {
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

	t.Run("RememberFunction", func(t *testing.T) {
		entity := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		value1, err := client.RememberForever(ctx, entity.ID, func(ctx context.Context) (*TestEntity, error) {
			return entity, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, entity.ID, value1.ID, "expected the value from the func to equal")

		value2, err := client.Remember(ctx, entity.ID, time.Minute*5, func(ctx context.Context) (*TestEntity, error) {
			return nil, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, entity.ID, value2.ID, "expected the value from the func to equal")
	})

	t.Run("ManyFunctions", func(t *testing.T) {
		key := "entity-list"

		entity1 := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		entity2 := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		entityList := []*TestEntity{entity1, entity2}

		err := client.PutManyForever(ctx, key, entityList)

		if err != nil {
			t.Error(err)
			return
		}

		has, err := client.Has(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.True(t, has)

		fetchedList, err := client.GetMany(ctx, key)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Len(t, fetchedList, 2, "expected fetched list to be length of 2")
		assert.Equal(t, entity1.ID, fetchedList[0].ID)
		assert.Equal(t, entity2.ID, fetchedList[1].ID)

	})

	t.Run("RememberManyFunctions", func(t *testing.T) {
		key := "remember-entity-list"

		entity1 := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		entity2 := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		entityList := []*TestEntity{entity1, entity2}

		fetchedList1, err := client.RememberMany(ctx, key, time.Minute*5, func(ctx context.Context) ([]*TestEntity, error) {
			return entityList, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Len(t, fetchedList1, 2, "expected fetched list to be length of 2")
		assert.Equal(t, entity1.ID, fetchedList1[0].ID)
		assert.Equal(t, entity2.ID, fetchedList1[1].ID)

		fetchedList2, err := client.RememberManyForever(ctx, key, func(ctx context.Context) ([]*TestEntity, error) {
			return []*TestEntity{}, nil
		})

		if err != nil {
			t.Error(err)
			return
		}

		assert.Len(t, fetchedList2, 2, "expected fetched list to be length of 2")
		assert.Equal(t, entity1.ID, fetchedList2[0].ID)
		assert.Equal(t, entity2.ID, fetchedList2[1].ID)

	})

	t.Run("ForgetFunctions", func(t *testing.T) {

		entity1 := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		err := client.Put(ctx, fmt.Sprintf("prefix-%s", entity1.ID), entity1, time.Minute*1)

		if err != nil {
			t.Error(err)
			return
		}

		entity2 := &TestEntity{
			ID:   gofakeit.UUID(),
			Name: gofakeit.Name(),
		}

		err = client.Put(ctx, fmt.Sprintf("prefix-%s", entity2.ID), entity2, time.Minute*1)

		if err != nil {
			t.Error(err)
			return
		}

		err = client.ForgetWithPrefix(ctx, "prefix-*")

		if err != nil {
			t.Error(err)
			return
		}

		has1, err := client.Has(ctx, fmt.Sprintf("prefix-%s", entity1.ID))

		if err != nil {
			t.Error(err)
			return
		}

		has2, err := client.Has(ctx, fmt.Sprintf("prefix-%s", entity2.ID))

		if err != nil {
			t.Error(err)
			return
		}

		assert.False(t, has1)
		assert.False(t, has2)

	})
}
