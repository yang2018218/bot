// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import "github.com/spf13/pflag"

type UtilsOptions struct {
	IP2LocationDBPath string `json:"ip2location_db_path"   mapstructure:"ip2location-db-path"`
}

func NewUtilsOptions() *UtilsOptions {
	return &UtilsOptions{}
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *UtilsOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&o.IP2LocationDBPath, "utils.ip2location_db_path", o.IP2LocationDBPath, "utils.ip2location_db_path")
}
