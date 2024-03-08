package app

import (
	"wechatbot/internal/pkg/core"
	metav1 "wechatbot/internal/wechatbot/meta/v1"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) Config(c *gin.Context) {
	data, err := ctrl.srv.App().Config(c)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	res := metav1.CommonResponse{
		Data: data,
	}
	core.WriteResponse(c, nil, res)
}
