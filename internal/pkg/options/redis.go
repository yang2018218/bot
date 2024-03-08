// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"wechatbot/pkg/db"

	"github.com/redis/go-redis/v9"
)

// RedisOptions defines options for mysql database.
type RedisOptions struct {
	Addr     string `json:"addr" mapstructure:"addr"`         // "localhost:6379",
	Password string `json:"password" mapstructure:"password"` // no password set
	DB       *int   `json:"db" mapstructure:"db"`             // use default DB
	Name     string `json:"name" mapstructure:"name"`
}

// NewRedisOptions create a `zero` value instance.
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{}
}

// NewClient create mysql store with the given config.
func (o *RedisOptions) NewClient() (*redis.Client, error) {
	opts := &db.RedisOptions{
		Addr:     o.Addr,
		Password: o.Password,
		DB:       o.DB,
	}

	return db.NewRedis(opts)
}
