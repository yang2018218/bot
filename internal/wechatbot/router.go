package wechatbot

import (
	appCtrlV1 "wechatbot/internal/wechatbot/controller/v1/app"
	"wechatbot/internal/wechatbot/httputils"
	"wechatbot/internal/wechatbot/store/postgres"

	"github.com/gin-gonic/gin"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}
func installController(g *gin.Engine) *gin.Engine {
	storeIns, _ := postgres.GetPostgreStore(nil)
	g.Use(httputils.ApiLogMiddleware(storeIns))
	r := g.Group("/api")
	appV1 := appCtrlV1.New(storeIns)
	appV1.RouterRegister(r.Group("/v1/app"))
	return g
}
