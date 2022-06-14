package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	responsepb "github.com/go-goim/api/transport/response"
	friendpb "github.com/go-goim/api/user/friend/v1"
	cgrpc "github.com/go-goim/core/pkg/conn/grpc"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/initialize"

	"github.com/go-goim/core/pkg/log"

	messagev1 "github.com/go-goim/api/message/v1"

	"github.com/go-goim/core/pkg/mq"

	"github.com/go-goim/gateway/internal/app"
)

type SendMessageService struct {
	messagev1.UnimplementedSendMessagerServer
	cp *cgrpc.ConnPool
}

var (
	sendMessageService = &SendMessageService{}
)

func GetSendMessageService() *SendMessageService {
	return sendMessageService
}

func init() {
	initialize.Register(initialize.NewBasicInitializer("send_message_service", nil, func() error {
		return sendMessageService.initConnPool()
	}))
}

func (s *SendMessageService) SendMessage(ctx context.Context, req *messagev1.SendMessageReq) (*messagev1.SendMessageResp, error) {
	rsp := new(messagev1.SendMessageResp)

	// check is friend
	if err := s.checkCanSendMsg(ctx, req); err != nil {
		rsp.Response = responsepb.NewBaseResponseWithMessage(responsepb.Code_RelationNotExist, err.Error())
		return nil, rsp.Response
	}

	mm := &messagev1.MqMessage{
		FromUser:        req.GetFromUser(),
		ToUser:          req.GetToUser(),
		PushMessageType: messagev1.PushMessageType_User,
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

func (s *SendMessageService) checkCanSendMsg(ctx context.Context, req *messagev1.SendMessageReq) error {
	cc, err := s.cp.Get()
	if err != nil {
		return err
	}

	resp, err := friendpb.NewFriendServiceClient(cc).IsFriend(ctx, &friendpb.BaseFriendRequest{
		Uid:       req.GetFromUser(),
		FriendUid: req.GetToUser(),
	})
	if err != nil {
		return err
	}

	if !resp.Success() {
		return resp
	}

	return nil
}

func (s *SendMessageService) initConnPool() error {
	cp, err := cgrpc.NewConnPool(cgrpc.WithInsecure(),
		cgrpc.WithClientOption(
			grpc.WithEndpoint(fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.UserService)),
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
