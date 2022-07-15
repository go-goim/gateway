package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	messagev1 "github.com/go-goim/api/message/v1"
	friendv1 "github.com/go-goim/api/user/friend/v1"
	sessionv1 "github.com/go-goim/api/user/session/v1"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/mq"
	"github.com/go-goim/core/pkg/util/snowflake"
	"github.com/go-goim/gateway/internal/dto"

	"github.com/go-goim/gateway/internal/app"
)

type SendMessageService struct{}

var (
	sendMessageService = &SendMessageService{}
)

func GetSendMessageService() *SendMessageService {
	return sendMessageService
}

func (s *SendMessageService) SendMessage(ctx context.Context, req *dto.SendMessageReq) (*dto.SendMessageResp, error) {
	pbReq := req.ToPb()
	// check is friend
	sid, err := s.checkCanSendMsg(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	mm := &messagev1.Message{
		From:        pbReq.From,
		To:          pbReq.To,
		SessionType: pbReq.SessionType,
		ContentType: pbReq.ContentType,
		Content:     pbReq.Content,
		SessionId:   sid,
		MsgId:       snowflake.Generate().Int64(),
		CreateTime:  time.Now().UnixMilli(),
	}

	err = s.sendMessage(ctx, mm)
	if err != nil {
		return nil, err
	}

	return &dto.SendMessageResp{
		SessionID: sid,
		MessageID: mm.MsgId,
	}, nil
}

func (s *SendMessageService) checkCanSendMsg(ctx context.Context, req *messagev1.SendMessageReq) (string, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return "", err
	}

	cr := &friendv1.CheckSendMessageAbilityRequest{
		FromUid:     req.From,
		ToUid:       req.To,
		SessionType: sessionv1.SessionType(req.SessionType),
	}

	resp, err := friendv1.NewFriendServiceClient(cc).CheckSendMessageAbility(ctx, cr)
	if err != nil {
		return "", err
	}

	if !resp.Response.Success() {
		return "", resp.Response
	}

	if resp.SessionId == nil || *resp.SessionId == "" {
		return "", errors.New("session id is nil")
	}

	return *resp.SessionId, nil
}

func (s *SendMessageService) Broadcast(ctx context.Context, req *dto.SendMessageReq) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)
	pbReq := req.ToPb()

	mm := &messagev1.Message{
		MsgId: snowflake.Generate().Int64(),
		// TODO: need session id for broadcast
		From:        pbReq.From,
		To:          pbReq.To,
		SessionType: sessionv1.SessionType_Broadcast,
		ContentType: pbReq.ContentType,
		Content:     pbReq.Content,
		CreateTime:  time.Now().UnixMilli(),
	}

	err := s.sendMessage(ctx, mm)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SendMessageService) sendMessage(ctx context.Context, mm *messagev1.Message) error {
	b, err := json.Marshal(mm)
	if err != nil {
		return err
	}

	// todo: maybe use another topic for all broadcast messages
	rs, err := app.GetApplication().Producer.SendSync(ctx, mq.NewMessage("def_topic", b))
	if err != nil {
		return err
	}

	log.Info("send message success", "rs", rs)
	return nil
}
