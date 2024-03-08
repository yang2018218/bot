package postgres

import (
	"context"
	"fmt"
	"sync"
	"wechatbot/internal/pkg/code"
	genericoptions "wechatbot/internal/pkg/options"
	appCfg "wechatbot/internal/wechatbot/config"
	"wechatbot/internal/wechatbot/store"
	"wechatbot/pkg/db"
	"wechatbot/pkg/errors"
	"wechatbot/pkg/log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	postgresFactory store.IStore
	DB              *gorm.DB
	once            sync.Once
)

type transactionKey struct{}

type datastore struct {
	db    *gorm.DB
	cache *redis.Client
	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *datastore) Cache() *redis.Client {
	return ds.cache
}

func (ds *datastore) SetCache(rc *redis.Client) {
	ds.cache = rc
}

func (ds *datastore) DB() *gorm.DB {
	return ds.db
}

func (ds *datastore) Core(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if ok {
		return tx
	}

	return ds.db
}

func (ds *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}
	log.Info("pg closed")
	return db.Close()
}

// GetFactoryOr create factory with the given config.
func GetPostgreStore(opts *genericoptions.PostgresOptions) (store.IStore, error) {
	if opts == nil && postgresFactory == nil {
		return nil, fmt.Errorf("failed to get postgresql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.PostgresOptions{
			DSN:      opts.DSN,
			Host:     opts.Host,
			Port:     opts.Port,
			User:     opts.User,
			Password: opts.Password,
			DBName:   opts.DBName,
			LogLevel: opts.LogLevel,
			// Logger:   logger.New(opts.LogLevel),
		}
		dbIns, err = db.NewPostgres(options)
		if appCfg.GCfg.GenericServerRunOptions.Mode == code.ServerModelDebug {
			dbIns = dbIns.Debug()
		}
		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)

		postgresFactory = &datastore{
			db: dbIns,
		}
	})

	if postgresFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get postgresql store fatory: %+v, error: %w", postgresFactory, err)
	}

	return postgresFactory, nil
}
