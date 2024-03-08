package service

import (
	"os"
	"wechatbot/internal/pkg/code"
	"wechatbot/internal/pkg/core"
	"wechatbot/internal/pkg/utils/netutil"
	"wechatbot/internal/wechatbot/config"
	"wechatbot/internal/wechatbot/store"
	"wechatbot/pkg/errors"
	"wechatbot/pkg/log"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Init() (err error)
	App() AppSrv
	Utils() UtilsSrv
	WechatBot() WechatBotSrv
}

type service struct {
	store store.IStore
}

func NewService(store store.IStore) Service {
	return &service{
		store: store,
	}
}

func (s *service) App() AppSrv {
	return newApp(s)
}

func (s *service) Utils() UtilsSrv {
	return newUtils(s)
}

func (s *service) WechatBot() WechatBotSrv {
	return newWechatBot(s)
}

var serviceInit bool

func (s *service) Init() (err error) {
	// 初始化IP2LocationDB
	if config.GCfg.UtilsOptions.IP2LocationDBPath != "" {
		if _, err := os.Stat(config.GCfg.UtilsOptions.IP2LocationDBPath); err == nil {
			_, err = netutil.OpenDB(config.GCfg.UtilsOptions.IP2LocationDBPath)
			if err != nil {
				log.Errorf("open ip2location db failed %s:%s", config.GCfg.UtilsOptions.IP2LocationDBPath, err.Error())
			} else {
				log.Infof("open ip2location db success %s", config.GCfg.UtilsOptions.IP2LocationDBPath)
			}
		}
	}
	log.Info("app init ok")
	serviceInit = true
	return
}

func Healthz(c *gin.Context) {
	if serviceInit {
		c.Set("log_level", -1)
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	} else {
		core.WriteResponse(c, errors.WithCode(code.ErrInternalServerError, "Not Ready Yet"), nil)
	}
}
