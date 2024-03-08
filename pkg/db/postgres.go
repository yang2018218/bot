// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"wechatbot/pkg/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Options defines optsions for mysql database.
type PostgresOptions struct {
	DSN      string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	LogLevel int
	Logger   logger.Interface
}

// New create a new gorm db instance with the given options.
func NewPostgres(opts *PostgresOptions) (*gorm.DB, error) {
	dsn := opts.DSN
	if dsn == "" {
		dsn = fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s`,
			opts.Host,
			opts.User,
			opts.Password,
			opts.DBName,
			opts.Port,
			"disable",
			"Asia/Shanghai",
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: opts.Logger,
		// 默认不启动事务
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	log.Infof("connect postgresql using dsn:%s success", dsn)

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	return nil, err
	// }

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}
