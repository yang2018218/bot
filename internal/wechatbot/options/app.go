// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/spf13/pflag"
)

type AppOptions struct {
	Version    string `json:"version"       mapstructure:"version"`
	AssetsPath string `json:"assets_path"    mapstructure:"assets-path"`
	Storage    string `json:"storage" mapstructure:"storage"`
	Setu       string `json:"setu" mapstructure:"setu"`
}

func NewAppOptions() *AppOptions {
	return &AppOptions{}
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *AppOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&o.Version, "app.version", o.Version, "app.version")
	fs.StringVar(&o.AssetsPath, "app.assets_path", o.Version, "app.assets_path")
	fs.StringVar(&o.Storage, "flyway", o.Storage, "storage")
}
