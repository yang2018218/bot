package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"wechatbot/internal/pkg/utils"
	"wechatbot/internal/wechatbot/define"
	metav1 "wechatbot/internal/wechatbot/meta/v1"
	"wechatbot/internal/wechatbot/model"
	"wechatbot/internal/wechatbot/store"

	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid"
)

const (
	// XRequestIDKey defines X-Request-ID key string.
	XRequestIDKey = "X-Request-ID"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ApiLogMiddleware(store store.IStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rid := c.GetHeader(XRequestIDKey)

		if rid == "" {
			rid = uuid.Must(uuid.NewV4()).String()
			c.Request.Header.Set(XRequestIDKey, rid)
			c.Set(XRequestIDKey, rid)
		}

		// Set XRequestIDKey header
		c.Writer.Header().Set(XRequestIDKey, rid)
		var (
			postRawByte []byte
			err         error
			blw         *bodyLogWriter
		)
		startTime := time.Now()
		rmeta := metav1.RequestMeta{
			Now: startTime,
		}
		c.Copy().ShouldBindHeader(&rmeta)
		rmeta.Host = c.Request.Host
		apiLog := model.ApiLogModel{
			Mehtod:   c.Request.Method,
			Type:     define.ApiLogTypeReceive,
			Path:     c.Request.URL.Path,
			RawQuery: c.Request.URL.RawQuery,
		}
		requestHeader := metav1.RequestHeader{}
		c.Copy().ShouldBindHeader(&requestHeader)
		// fmt.Printf("middleware log \n")
		// for name, values := range c.Request.Header {
		// 	valueString := strings.Join(values, ";")
		// 	fmt.Printf("%s:%s \n", name, valueString)
		// }
		headers, _ := utils.Struct2Map(&requestHeader)
		apiLog.RequestHeader = headers
		if apiLog.Mehtod == "POST" {
			postRawByte, err = c.GetRawData()
			rmeta.PostRawByte = postRawByte
			if err != nil {
				fmt.Print(err)
			}
			if strings.Contains(rmeta.ContentType, "application/x-www-form-urlencoded") || strings.Contains(rmeta.ContentType, "application/form-data") {
				apiLog.RawRequestBody = string(postRawByte)
			} else {
				var requestBody interface{}
				err = json.Unmarshal(postRawByte, &requestBody)
				if err == nil {
					apiLog.RequestBody = requestBody
				} else {
					apiLog.RequestBody = map[string]string{
						"raw": string(postRawByte),
					}
				}
			}

			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(postRawByte))
			blw = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}
		if rmeta.ClientIP == "" {
			rmeta.ClientIP = c.ClientIP()
		}
		operatorEmpty := true
		operator := metav1.Operator{}
		for name, values := range c.Request.Header {
			lname := strings.ToLower(name)
			valueString := strings.Join(values, ";")
			if valueString == "" {
				continue
			}
			if utils.ContainsString([]string{"openid", "x-openid"}, lname) {
				rmeta.OpenID = valueString
				operator.OpenID = rmeta.OpenID
				operatorEmpty = false
			} else if utils.ContainsString([]string{"unionid", "x-unionid"}, lname) {
				rmeta.UnionID = valueString
				operator.UnionID = rmeta.UnionID
				operatorEmpty = false
			} else if utils.ContainsString([]string{"appid", "x-appid"}, lname) {
				rmeta.AppID = valueString
				operator.AppID = rmeta.AppID
				operatorEmpty = false
			} else if utils.ContainsString([]string{"accountid", "x-accountid"}, lname) {
				operator.AccountID = valueString
				operatorEmpty = false
			} else if utils.ContainsString([]string{"account", "x-account"}, lname) {
				operator.Account = valueString
				operatorEmpty = false
			} else if utils.ContainsString([]string{"phone", "x-phone"}, lname) {
				operator.Phone = valueString
				operatorEmpty = false
			}
		}

		apiLog.Host = rmeta.Host
		apiLog.RequestID = rmeta.RequestID
		apiLog.IP = rmeta.ClientIP

		c.Set("rmeta", rmeta)
		if !operatorEmpty {
			apiLog.Operator = &operator
			c.Set("operator", operator)
		}
		// 设置网络信息
		c.Next()
		level := c.GetInt("log_level")
		apiLog.Cost = time.Since(startTime).Milliseconds()
		ra := time.Now()
		apiLog.ResponsedAt = &ra
		apiLog.Status = c.Writer.Status()
		if level == int(define.ApiLogLevel1) && blw != nil {
			var body interface{}
			err = json.Unmarshal(blw.body.Bytes(), &body)
			if err == nil {
				apiLog.ResponseBody = body
			} else {
				apiLog.ResponseBody = map[string]string{
					"raw": blw.body.String(),
				}
			}
		}
		if level > 0 {
			// save
			// go store.ApiLog().Create(c, &apiLog)
		}
	}
}
