package service

import (
	"context"
	"encoding/json"
	"errors"

	responsepb "github.com/go-goim/api/transport/response"
	friendpb "github.com/go-goim/api/user/friend/v1"
	sessionpb "github.com/go-goim/api/user/session/v1"
	"github.com/go-goim/core/pkg/log"

	messagev1 "github.com/go-goim/api/message/v1"

	"github.com/go-goim/core/pkg/mq"

	"github.com/go-goim/gateway/internal/app"
)

type SendMessageService struct {
	messagev1.UnimplementedSendMessagerServer
}

var (
	sendMessageService = &SendMessageService{}
)

func GetSendMessageService() *SendMessageService {
	return sendMessageService
}

func (s *SendMessageService) SendMessage(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)

	// check is friend
	sid, err := s.checkCanSendMsg(ctx, req)
	if err != nil {
		rsp.Response = responsepb.NewBaseResponseWithError(err)
		return nil, rsp.Response
	}

	mm := &messagev1.MqMessage{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_User,
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
		SessionId:       sid,
	}

	rsp, err = s.sendMessage(ctx, mm)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.Response
	}

	rsp.SessionId = sid
	return rsp, nil
}

func (s *SendMessageService) checkCanSendMsg(ctx context.Context, req *messagev1.SendMessageReq) (int64, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return 0, err
	}

	cr := &friendpb.CheckSendMessageAbilityRequest{
		FromUid:     req.GetFromUser(),
		ToUid:       req.GetToUser(),
		SessionType: sessionpb.SessionType_SingleChat,
	}

	// todo check touid whether is a group id

	resp, err := friendpb.NewFriendServiceClient(cc).CheckSendMessageAbility(ctx, cr)
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
		rsp.Response = responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error())
		return nil, rsp.Response
	}

	mm := &messagev1.MqMessage{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_Broadcast,
		ContentType:     req.GetContentType(),
		Content:         req.GetContent(),
	}

	rsp, err := s.sendMessage(ctx, mm)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.Response
	}

	return rsp, nil
}

func (s *SendMessageService) sendMessage(ctx context.Context, mm *messagev1.MqMessage) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)
	rsp.Response = responsepb.Code_OK.BaseResponse()

	b, err := json.Marshal(mm)
	if err != nil {
		rsp.Response = responsepb.NewBaseResponseWithError(err)
		return rsp, nil
	}

	// todo: maybe use another topic for all broadcast messages
	rs, err := app.GetApplication().Producer.SendSync(ctx, mq.NewMessage("def_topic", b))
	if err != nil {
		rsp.Response = responsepb.NewBaseResponseWithError(err)
		return rsp, nil
	}

	log.Info("send message success", "rs", rs)
	rsp.MsgSeq = rs.MsgID

	return rsp, nil
}
