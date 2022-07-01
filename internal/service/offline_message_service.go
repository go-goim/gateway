package service

import (
	"context"

	messagev1 "github.com/go-goim/api/message/v1"
)

type OfflineMessageService struct {
}

var (
	offlineMsgSrc = &OfflineMessageService{}
)

func GetOfflineMessageService() *OfflineMessageService {
	return offlineMsgSrc
}

func (s *OfflineMessageService) QueryOfflineMsg(ctx context.Context, req *messagev1.QueryOfflineMessageReq) (
	[]*messagev1.BriefMessage, error) {
	cc, err := msgServiceConnPool.Get()
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
