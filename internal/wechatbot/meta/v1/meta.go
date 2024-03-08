package v1

import (
	"context"
	"time"
	"wechatbot/internal/wechatbot/define"
)

type CommonResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ListResponse struct {
	CommonResponse `json:",inline"`
	Count          int64  `json:"count"`
	HasMore        bool   `json:"has_more"`
	ExportUrl      string `json:"export_url,omitempty"`
}

type ListParams struct {
	Offset     int `json:"offset" form:"offset"`
	Limit      int `json:"limit" form:"limit"`
	OrderIndex int `json:"order_index" form:"order_index"`
	OrderValue int `json:"order_value" form:"order_value"`
}

type ListMeta struct {
	Count   int64
	HasMore bool
}

type RequestHeader struct {
	UserAgent      string `header:"User-Agent" json:"user_agent,omitempty"`
	ContentType    string `header:"Content-Type" json:"content_type,omitempty"`
	AcceptLanguage string `header:"Accept-Language" json:"accept_language,omitempty"`
	Cookie         string `header:"Cookie" json:"-"`
	Origin         string `header:"Origin" json:"origin,omitempty"`
	Referer        string `header:"Referer"  json:"referer,omitempty"`
	XRealIP        string `header:"X-Real-IP" json:"-"`
	XForwardedFor  string `header:"X-Forwarded-For" json:"-"`
}

type RequestMeta struct {
	RequestHeader
	Now         time.Time
	ClientIP    string `header:"X-Real-Ip"`
	App         string `header:"App"`
	AppEnv      string `header:"App-Env"`
	AppID       string `header:"Appid"`
	OpenID      string `header:"Openid"`
	UnionID     string `header:"Unionid"`
	RequestID   string `header:"X-Request-ID"`
	Host        string
	ClientType  define.ClientType
	PostRawByte []byte
}

type Operator struct {
	AccountID interface{} `json:"account_id,omitempty"`
	ID        interface{} `json:"id,omitempty"`
	Phone     string      `json:"phone,omitempty"`
	Name      string      `json:"name,omitempty"`
	Account   string      `json:"account,omitempty"`
	OpenID    string      `json:"openid,omitempty"`
	UnionID   string      `json:"unionid,omitempty"`
	AppID     string      `json:"appid,omitempty"`
}

type GetDetailParams struct {
	ID uint64 `form:"id" binding:"required"`
}

type Address struct {
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
	Address  string `json:"address,omitempty"`
}

type DeleteParams struct {
	ID uint64 `json:"id" binding:"required"`
}

func GetContextInfo(ctx context.Context) (rmeta RequestMeta, operator *Operator) {
	rmetai := ctx.Value("rmeta")
	if rmetai != nil {
		rmeta = rmetai.(RequestMeta)
	}
	if rmeta.Now.IsZero() {
		rmeta.Now = time.Now()
	}
	operatorI := ctx.Value("operator")
	if operatorI != nil {
		switch operatorI.(type) {
		case *Operator:
			operator = operatorI.(*Operator)
		}
	}
	return
}
