package service

import (
	"context"
	"encoding/json"
	"errors"

	messagev1 "github.com/go-goim/api/message/v1"
	friendv1 "github.com/go-goim/api/user/friend/v1"
	sessionv1 "github.com/go-goim/api/user/session/v1"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/mq"
	"github.com/go-goim/core/pkg/util"
	"github.com/go-goim/core/pkg/util/snowflake"

	"github.com/go-goim/gateway/internal/app"
)

type SendMessageService struct{}

var (
	sendMessageService = &SendMessageService{}
)

func GetSendMessageService() *SendMessageService {
	return sendMessageService
}

func (s *SendMessageService) SendMessage(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	// check is friend
	sid, err := s.checkCanSendMsg(ctx, req)
	if err != nil {
		return nil, err
	}

	mm := &messagev1.Message{
		From:        req.GetFrom(),
		To:          req.GetTo(),
		SessionType: sessionv1.SessionType_SingleChat,
		ContentType: req.GetContentType(),
		Content:     req.GetContent(),
		SessionId:   sid,
		MsgId:       snowflake.Generate().Int64(),
	}

	if util.IsGroupUID(req.GetTo()) {
		mm.SessionType = sessionv1.SessionType_GroupChat
	}

	err = s.sendMessage(ctx, mm)
	if err != nil {
		return nil, err
	}

	return &messagev1.SendMessageResp{
		SessionId: sid,
		MsgId:     mm.MsgId,
	}, nil
}

func (s *SendMessageService) checkCanSendMsg(ctx context.Context, req *messagev1.SendMessageReq) (int64, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return 0, err
	}

	cr := &friendv1.CheckSendMessageAbilityRequest{
		FromUid:     req.GetFrom(),
		ToUid:       req.GetTo(),
		SessionType: sessionv1.SessionType_SingleChat,
	}

	if util.IsGroupUID(req.GetTo()) {
		cr.SessionType = sessionv1.SessionType_GroupChat
	}

	resp, err := friendv1.NewFriendServiceClient(cc).CheckSendMessageAbility(ctx, cr)
	if err != nil {
		return 0, err
	}

	if !resp.Response.Success() {
		return 0, resp.Response
	}

	if resp.SessionId == nil || *resp.SessionId == 0 {
		return 0, errors.New("session id is nil")
	}

	return *resp.SessionId, nil
}

func (s *SendMessageService) Broadcast(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)
	// check req params
	if err := req.Validate(); err != nil {
		return nil, err
	}

	mm := &messagev1.Message{
		MsgId: snowflake.Generate().Int64(),
		// TODO: need session id for broadcast
		From:        req.GetFrom(),
		To:          req.GetTo(),
		SessionType: sessionv1.SessionType_Broadcast,
		ContentType: req.GetContentType(),
		Content:     req.GetContent(),
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
