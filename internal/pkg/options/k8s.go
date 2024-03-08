// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"time"

	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type K8SOptions struct {
	Namespace   string                `json:"namespace"       mapstructure:"namespace"`
	PodName     string                `json:"pod_name"        mapstructure:"pod-name"`
	ServiceHost string                `json:"service_host"`
	LeaseName   string                `json:"lease_name"      mapstructure:"lease-name"`
	IsLeader    bool                  `json:"is_leader"`
	Clientset   *kubernetes.Clientset `json:"-"`
	LeaseTime   []time.Time           `json:"lease_time"`
}

func NewK8SOptions() *K8SOptions {
	return &K8SOptions{}
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *K8SOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	// fs.StringVar(&o.Key, "sms.key", o.Key, "sms.key")
	// fs.StringVar(&o.Secret, "sms.secret", o.Secret, "sms.secret")
	// fs.StringVar(&o.Endpoint, "sms.endpoint", o.Endpoint, "sms.endpoint")
	// fs.StringVar(&o.ApiVersion, "sms.apiversion", o.ApiVersion, "sms.apiversion")
}
