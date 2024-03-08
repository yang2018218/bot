package app

import (
	srv "wechatbot/internal/wechatbot/service"
	"wechatbot/internal/wechatbot/store"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	srv srv.Service
}

// NewUserController creates a user handler.
func New(store store.IStore) *Controller {
	return &Controller{
		srv: srv.NewService(store),
	}
}

func (ctrl *Controller) RouterRegister(g *gin.RouterGroup) {
	g.GET("/env", ctrl.Env)
	g.GET("/config", ctrl.Config)
}
