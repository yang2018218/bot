package service

import (
	"context"
	"os"
	"strings"
	"wechatbot/internal/wechatbot/config"
	"wechatbot/internal/wechatbot/store"
)

type AppSrv interface {
	Env(ctx context.Context) (map[string]string, error)
	Config(ctx context.Context) (interface{}, error)
}

type appService struct {
	store store.IStore
}

var _ AppSrv = (*appService)(nil)

func newApp(srv *service) *appService {
	return &appService{
		store: srv.store,
	}
}

func (s *appService) Env(ctx context.Context) (map[string]string, error) {
	envMap := make(map[string]string)
	// 获取全部的环境变量
	allenv := os.Environ()

	// 打印环境变量
	for _, e := range allenv {
		result := strings.SplitN(e, "=", 2)
		envMap[result[0]] = result[1]
	}
	return envMap, nil
}

func (s *appService) Config(ctx context.Context) (interface{}, error) {
	return config.GCfg, nil
}
