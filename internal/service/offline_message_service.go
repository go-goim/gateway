package service

import (
	"context"

	messagev1 "github.com/go-goim/api/message/v1"
	"github.com/go-goim/gateway/internal/dto"
)

type OfflineMessageService struct {
}

var (
	offlineMsgSrc = &OfflineMessageService{}
)

func GetOfflineMessageService() *OfflineMessageService {
	return offlineMsgSrc
}

func (s *OfflineMessageService) QueryOfflineMsg(ctx context.Context, req *dto.QueryOfflineMessageReq) (
	[]*dto.Message, int32, error) {
	cc, err := msgServiceConnPool.Get()
	if err != nil {
		return nil, 0, err
	}

	rsp, err := messagev1.NewOfflineMessageServiceClient(cc).QueryOfflineMessage(ctx, req.ToPb())
	if err != nil {
		return nil, 0, err
	}

	if err := rsp.GetError().Err(); err != nil {
		return nil, 0, err
	}

	return dto.MessagesFromPb(rsp.Messages), rsp.GetTotal(), nil
}
