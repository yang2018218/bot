package app

import (
	"wechatbot/internal/pkg/core"
	metav1 "wechatbot/internal/wechatbot/meta/v1"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) Env(c *gin.Context) {
	data, err := ctrl.srv.App().Env(c)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	res := metav1.CommonResponse{
		Data: data,
	}
	core.WriteResponse(c, nil, res)
}
