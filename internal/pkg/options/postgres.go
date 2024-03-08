// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/spf13/pflag"
	"gorm.io/gorm"

	"wechatbot/pkg/db"
)

// PostgresOptions defines options for mysql database.
type PostgresOptions struct {
	DSN        string `json:"dsn,omitempty"                     mapstructure:"dsn"`
	Host       string `json:"host,omitempty"                     mapstructure:"host"`
	Port       int    `json:"port,omitempty"                 mapstructure:"port"`
	User       string `json:"user,omitempty"                 mapstructure:"user"`
	Password   string `json:"-"                                  mapstructure:"password"`
	DBName     string `json:"dbname"                                  mapstructure:"dbname"`
	LogLevel   int    `json:"log-level"                          mapstructure:"log-level"`
	FlywayPath string `json:"flyway-path" mapstructure:"flyway-path"`
}

// NewPostgresOptions create a `zero` value instance.
func NewPostgresOptions() *PostgresOptions {
	return &PostgresOptions{
		DSN:        "",
		Host:       "127.0.0.1",
		Port:       5432,
		User:       "postgres",
		Password:   "",
		DBName:     "postgres",
		LogLevel:   1, // Silent
		FlywayPath: "",
	}
}

// Validate verifies flags passed to PostgresOptions.
func (o *PostgresOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *PostgresOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.DSN, "postgres.dsn", o.DSN, ""+
		"PostgreSQL service host address. If not blank, the following related options will be ignored.")

	fs.StringVar(&o.Host, "postgres.host", o.Host, ""+
		"PostgreSQL service host address. If left blank and dsn is blank, the following related mysql options will be ignored.")

	fs.StringVar(&o.User, "postgres.user", o.User, ""+
		"User for access to postgres service.")

	fs.StringVar(&o.Password, "postgres.password", o.Password, ""+
		"Password for access to postgres, should be used pair with password.")

	fs.StringVar(&o.DBName, "postgres.dbname", o.DBName, ""+
		"Database name for the server to use.")

	fs.IntVar(&o.LogLevel, "postgres.log-level", o.LogLevel, ""+
		"Specify gorm log level.")

	fs.StringVar(&o.FlywayPath, "postgres.flyway-path", o.FlywayPath, ""+
		"Flyway path")
}

// NewClient create mysql store with the given config.
func (o *PostgresOptions) NewClient() (*gorm.DB, error) {
	opts := &db.PostgresOptions{
		DSN:      o.DSN,
		Host:     o.Host,
		Port:     o.Port,
		User:     o.User,
		Password: o.Password,
		DBName:   o.DBName,
		LogLevel: o.LogLevel,
	}

	return db.NewPostgres(opts)
}
