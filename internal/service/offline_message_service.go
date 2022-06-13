package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	messagev1 "github.com/go-goim/api/message/v1"
	cgrpc "github.com/go-goim/core/pkg/conn/grpc"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/initialize"

	"github.com/go-goim/gateway/internal/app"
)

type OfflineMessageService struct {
	cp *cgrpc.ConnPool
}

var (
	offlineMsgSrc = &OfflineMessageService{}
)

func GetOfflineMessageService() *OfflineMessageService {
	return offlineMsgSrc
}

func init() {
	initialize.Register(initialize.NewBasicInitializer("offline_message_service", nil, func() error {
		return offlineMsgSrc.initConnPool()
	}))
}

func (s *OfflineMessageService) QueryOfflineMsg(ctx context.Context, req *messagev1.QueryOfflineMessageReq) (
	[]*messagev1.BriefMessage, error) {
	cc, err := s.cp.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := messagev1.NewOfflineMessageClient(cc).QueryOfflineMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.Response
	}

	return rsp.GetMessages(), nil
}

func (s *OfflineMessageService) initConnPool() error {
	cp, err := cgrpc.NewConnPool(cgrpc.WithInsecure(),
		cgrpc.WithClientOption(
			grpc.WithEndpoint(fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.MsgService)),
			grpc.WithDiscovery(app.GetApplication().Register),
			grpc.WithTimeout(time.Second*5),
			grpc.WithOptions(ggrpc.WithBlock()),
		), cgrpc.WithPoolSize(2))
	if err != nil {
		return err
	}

	s.cp = cp
	graceful.Register(func(_ context.Context) error {
		return cp.Release()
	})
	return nil
}
