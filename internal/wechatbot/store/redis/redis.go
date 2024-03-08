package redis

import (
	"fmt"
	"sync"
	genericoptions "wechatbot/internal/pkg/options"
	"wechatbot/pkg/db"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// GetFactoryOr create factory with the given config.
func GetRedisCache(opts *genericoptions.RedisOptions) (*redis.Client, error) {
	if opts == nil && client == nil {
		return nil, fmt.Errorf("failed to get redis cache fatory")
	}

	var err error
	once.Do(func() {
		options := &db.RedisOptions{
			Addr:     opts.Addr,
			Password: opts.Password,
			DB:       opts.DB,
		}
		client, err = db.NewRedis(options)
		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)

		// postgresFactory = &datastore{dbIns}
	})

	if client == nil || err != nil {
		return nil, fmt.Errorf("failed to get redis cache fatory: %+v, error: %w", client, err)
	}

	return client, nil
}
