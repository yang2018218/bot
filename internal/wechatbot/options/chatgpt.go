// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/spf13/pflag"
)

type ChatGptOptions struct {
	ApiUrl string `json:"apiurl" mapstructure:"apiurl"`
	Key    string `json:"key"    mapstructure:"key"`
}

func NewChatGptOptions() *ChatGptOptions {
	return &ChatGptOptions{}
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *ChatGptOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&o.ApiUrl, "chatgpt.apiurl", o.ApiUrl, "chatgpt.apiurl")
	fs.StringVar(&o.Key, "chatgpt.key", o.Key, "chatgpt.key")
}
