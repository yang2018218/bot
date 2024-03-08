package httputils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wechatbot/internal/wechatbot/define"
	"wechatbot/internal/wechatbot/model"
	"wechatbot/pkg/errors"

	"gorm.io/gorm"
)

type RequestOption struct {
	RequestBody    interface{}
	RequestHeaders map[string]string
	RequestForm    *url.Values
	RespBody       interface{}
	ContentType    string
	SaveLog        bool
	LogDB          *gorm.DB
}

func Request(method string, _url string, opt *RequestOption) (err error) {
	u, errU := url.Parse(_url)
	if errU != nil {
		err = errors.Wrap(errU, "invalid url")
		return
	}
	var (
		req  *http.Request
		resp *http.Response
	)
	apiLog := model.ApiLogModel{
		Type:     define.ApiLogTypeRequest,
		Mehtod:   method,
		Host:     u.Host,
		Path:     u.Path,
		RawQuery: u.RawQuery,
	}
	defer func() {
		if opt != nil && opt.SaveLog && opt.LogDB != nil && req != nil {
			// apiLog.Create(opt.LogDB)
		}
	}()
	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, _url, nil)
		if err != nil {
			return
		}
		if opt != nil && opt.ContentType != "" {
			req.Header.Set("Content-Type", opt.ContentType)
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
	} else if method == http.MethodPost {
		if opt != nil && (opt.ContentType == "application/x-www-form-urlencoded;charset=utf-8" || opt.RequestForm != nil) {
			// form 上传
			encodedData := ""
			if opt.RequestForm != nil {
				encodedData = opt.RequestForm.Encode()
			}
			req, err = http.NewRequest(http.MethodPost, _url, strings.NewReader(encodedData))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
			req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
		} else {
			// json 上传
			var body = &bytes.Buffer{}
			if opt != nil && opt.RequestBody != nil {
				apiLog.RequestBody = opt.RequestBody
				json.NewEncoder(body).Encode(opt.RequestBody)
			}
			req, err = http.NewRequest(http.MethodPost, _url, body)
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		err = errors.New("error request")
		return
	}
	if opt != nil && len(opt.RequestHeaders) > 0 {
		for key, value := range opt.RequestHeaders {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	if len(req.Header) > 0 {
		RequestHeader := make(map[string]string)
		for k, v := range req.Header {
			if len(v) == 0 {
				continue
			}
			RequestHeader[k] = strings.Join(v, ";")
		}
		apiLog.RequestHeader = RequestHeader
	}
	apiLog.CreatedAt = time.Now()
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	rt := time.Now()
	apiLog.ResponsedAt = &rt
	apiLog.Cost = rt.Sub(apiLog.CreatedAt).Milliseconds()
	apiLog.Status = resp.StatusCode
	if len(resp.Header) > 0 {
		respHeader := make(map[string]string)
		for k, v := range resp.Header {
			if len(v) == 0 {
				continue
			}
			respHeader[k] = strings.Join(v, ";")
		}
		apiLog.ResponseHeader = respHeader
	}
	if opt != nil && opt.RespBody != nil {
		defer resp.Body.Close()
		respContentType := resp.Header.Get("Content-Type")
		if strings.Contains(respContentType, "application/json") {
			err = json.NewDecoder(resp.Body).Decode(opt.RespBody)
			if err != nil {
				return
			}
			apiLog.ResponseBody = opt.RespBody
		} else if strings.Contains(respContentType, "text/plain") {
			respbody, errResp := io.ReadAll(resp.Body)
			if errResp != nil {
				err = errResp
				return
			}
			bodyStr := string(respbody)
			opt.RespBody = &bodyStr
		}
	}
	return
}
