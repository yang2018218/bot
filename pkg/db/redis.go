// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package db

import (
	"github.com/redis/go-redis/v9"
)

// Options defines optsions for mysql database.
type RedisOptions struct {
	Addr     string
	Password string
	DB       *int
}

// New create a new gorm db instance with the given options.
func NewRedis(opts *RedisOptions) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password, // no password set
		DB:       *opts.DB,      // use default DB
	})
	return rdb, nil
}
