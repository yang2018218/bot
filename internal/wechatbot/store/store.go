package store

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var client IStore

type IStore interface {
	SetCache(*redis.Client)
	Cache() *redis.Client
	DB() *gorm.DB
	TX(context.Context, func(ctx context.Context) error) error
	Close() error
}

func Client() IStore {
	return client
}

func SetClient(factory IStore) {
	client = factory
}
