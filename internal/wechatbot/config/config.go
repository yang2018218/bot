// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"wechatbot/internal/wechatbot/options"

	"github.com/gin-gonic/gin"
)

// Config is the running configuration structure of the IAM pump service.
var (
	GCfg Config
)

type Config struct {
	*options.Options
}

type AuthMiddleware func(bool) gin.HandlerFunc

// CreateConfigFromOptions creates a running configuration instance based
// on a given IAM pump command line or configuration file option.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
