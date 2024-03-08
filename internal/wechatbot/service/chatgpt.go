package service

import (
	"context"
	"fmt"
	"net/http"
	"wechatbot/internal/wechatbot/config"
	"wechatbot/internal/wechatbot/httputils"
	metav1 "wechatbot/internal/wechatbot/meta/v1"
	"wechatbot/internal/wechatbot/store"
)

type ChatGptSrv interface {
	ChatCompletions(ctx context.Context, params metav1.ChatGPTChatCompletionsOpts) (*metav1.ChatGPTChatCompletionsResponse, error)
}

type chatGptService struct {
	store store.IStore
}

var _ ChatGptSrv = (*chatGptService)(nil)

func newChatGpt(srv *service) *chatGptService {
	return &chatGptService{
		store: srv.store,
	}
}

func (s *chatGptService) ChatCompletions(ctx context.Context, params metav1.ChatGPTChatCompletionsOpts) (*metav1.ChatGPTChatCompletionsResponse, error) {
	gptRes := metav1.ChatGPTChatCompletionsResponse{}
	requestOpt := httputils.RequestOption{
		RespBody:    &gptRes,
		RequestBody: params,
		RequestHeaders: map[string]string{
			"Authorization": "Bearer " + config.GCfg.ChatGptOptions.Key,
		},
	}
	err := httputils.Request(http.MethodPost, fmt.Sprintf("%s/v1/chat/completions", config.GCfg.ChatGptOptions.ApiUrl), &requestOpt)
	if err != nil {
		return nil, err
	}
	return &gptRes, nil
}
