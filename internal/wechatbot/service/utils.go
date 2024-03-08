package service

import (
	"context"
	"time"
	"wechatbot/internal/pkg/code"
	"wechatbot/internal/pkg/utils/netutil"
	metav1 "wechatbot/internal/wechatbot/meta/v1"
	"wechatbot/internal/wechatbot/store"
	"wechatbot/pkg/errors"
)

type UtilsSrv interface {
	IP2Location(ctx context.Context, params metav1.UtilsIP2LocationParams) (data interface{}, err error)
	Panic(ctx context.Context)
}

type utilsService struct {
	store store.IStore
}

var _ UtilsSrv = (*utilsService)(nil)

func newUtils(srv *service) *utilsService {
	return &utilsService{
		store: srv.store,
	}
}

func (s *utilsService) IP2Location(ctx context.Context, params metav1.UtilsIP2LocationParams) (data interface{}, err error) {
	if params.IP == "" {
		rmeta := ctx.Value("rmeta").(metav1.RequestMeta)
		params.IP = rmeta.ClientIP
	}
	if params.IP == "" {
		err = errors.WithCode(code.ErrInvalidParams, "no ip")
		return
	}
	data, err = netutil.IPv4ToLocation(params.IP)
	return
}

func (s *utilsService) Panic(ctx context.Context) {
	panic(time.Now().Format(time.DateTime))
}
