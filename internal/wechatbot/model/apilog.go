package model

import (
	"encoding/json"
	"net/url"
	"time"

	"wechatbot/internal/pkg/utils"
	"wechatbot/internal/wechatbot/define"

	metav1 "wechatbot/internal/wechatbot/meta/v1"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ApiLogModel struct {
	ID                 int64             `gorm:"column:id"`
	CreatedAt          time.Time         `gorm:"column:created_at"`
	OperatorJSON       datatypes.JSON    `gorm:"column:operator;default:null;type:jsonb"`
	Operator           *metav1.Operator  `gorm:"-"`
	Host               string            `gorm:"column:host;default:null"`
	IP                 string            `gorm:"column:ip;default:null"`
	Mehtod             string            `gorm:"column:method"`
	Path               string            `gorm:"column:path;default:null"`
	QueryJSON          datatypes.JSON    `gorm:"column:query;default:null"`
	RawQuery           string            `gorm:"-"`
	RequestID          string            `gorm:"column:request_id;default:null"`
	RawRequestBody     string            `gorm:"-"`
	RequestBodyJSON    string            `gorm:"column:request_body;default:null"`
	RequestBody        interface{}       `gorm:"-"`
	RequestHeaderJSON  datatypes.JSON    `gorm:"column:request_header;default:null"`
	RequestHeader      map[string]string `gorm:"-"`
	ResponseBodyJSON   string            `gorm:"column:response_body;default:null"`
	ResponseBody       interface{}       `gorm:"-"`
	ResponseHeaderJSON datatypes.JSON    `gorm:"column:response_header;default:null"`
	ResponseHeader     map[string]string `gorm:"-"`
	Error              string            `gorm:"column:error;default:null"`
	Type               define.ApiLogType `gorm:"column:type"`
	ResponsedAt        *time.Time        `gorm:"column:responsed_at"`
	Cost               int64             `gorm:"column:cost;default:null"` // 毫秒耗时
	Status             int               `gorm:"status;default:null"`
}

func (ApiLogModel) TableName() string {
	return "api_log"
}

func (m *ApiLogModel) Create(tx *gorm.DB) (err error) {
	err = tx.Create(m).Error
	return err
}

func (m *ApiLogModel) Updates(tx *gorm.DB) (err error) {
	err = tx.Model(m).Updates(m).Error
	return err
}

func (i *ApiLogModel) BeforeSave(tx *gorm.DB) (err error) {
	if i.RawQuery != "" {
		u := url.URL{RawQuery: i.RawQuery}
		q := u.Query()
		m := map[string]interface{}{}
		for k, v := range q {
			m[k] = v[0]
		}
		jsonByte, _ := json.Marshal(m)
		i.QueryJSON = datatypes.JSON(jsonByte)
	}

	if i.Operator != nil {
		jsonByte, _ := json.Marshal(i.Operator)
		i.OperatorJSON = datatypes.JSON(jsonByte)
	}
	if !utils.IsNil(i.RequestBody) {
		jsonByte, _ := json.Marshal(i.RequestBody)
		i.RequestBodyJSON = string(jsonByte)
	} else if i.RawRequestBody != "" {
		u := url.URL{RawQuery: i.RawRequestBody}
		q := u.Query()
		m := map[string]interface{}{}
		for k, v := range q {
			m[k] = v[0]
		}
		jsonByte, _ := json.Marshal(m)
		i.RequestBodyJSON = string(jsonByte)
	}
	if !utils.IsNil(i.ResponseBody) {
		jsonByte, _ := json.Marshal(i.ResponseBody)
		i.ResponseBodyJSON = string(jsonByte)
	}
	if len(i.RequestHeader) > 0 {
		jsonByte, _ := json.Marshal(i.RequestHeader)
		i.RequestHeaderJSON = datatypes.JSON(jsonByte)
	}
	if len(i.ResponseHeader) > 0 {
		jsonByte, _ := json.Marshal(i.ResponseHeader)
		i.ResponseHeaderJSON = datatypes.JSON(jsonByte)
	}
	return
}

func (i *ApiLogModel) AfterFind(tx *gorm.DB) (err error) {
	if len(i.OperatorJSON) > 0 {
		json.Unmarshal(i.OperatorJSON, &i.Operator)
	}
	if len(i.RequestBodyJSON) > 0 {
		utils.Unmarshal(i.RequestBodyJSON, &i.RequestBody)
	}
	if len(i.ResponseBodyJSON) > 0 {
		utils.Unmarshal(i.ResponseBodyJSON, &i.ResponseBody)
	}
	if len(i.RequestHeaderJSON) > 0 {
		json.Unmarshal(i.RequestHeaderJSON, &i.RequestHeader)
	}
	if len(i.ResponseHeaderJSON) > 0 {
		json.Unmarshal(i.ResponseHeaderJSON, &i.ResponseHeader)
	}
	return
}
