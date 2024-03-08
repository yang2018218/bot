// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"net/http"

	"wechatbot/pkg/errors"
	"wechatbot/pkg/log"

	"github.com/gin-gonic/gin"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		errStr := err.Error()
		if errStr == "" {
			errStr = coder.String()
		}
		c.AbortWithStatusJSON(coder.HTTPStatus(), ErrResponse{
			Code:      coder.Code(),
			Message:   errStr,
			Reference: coder.Reference(),
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
