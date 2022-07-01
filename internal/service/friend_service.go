package service

import (
	"context"

	responsepb "github.com/go-goim/api/transport/response"
	friendpb "github.com/go-goim/api/user/friend/v1"
)

type FriendService struct{}

var (
	friendService = &FriendService{}
)

func GetFriendService() *FriendService {
	return friendService
}

/*
 * Friend Request Logic
 */

func (s *FriendService) AddFriend(ctx context.Context, req *friendpb.AddFriendRequest) (*friendpb.AddFriendResult, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).AddFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetResult(), nil
}

func (s *FriendService) QueryFriendRequestList(ctx context.Context, req *friendpb.QueryFriendRequestListRequest) (
	[]*friendpb.FriendRequest, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).QueryFriendRequestList(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetFriendRequestList(), nil
}

func (s *FriendService) AcceptFriend(ctx context.Context, req *friendpb.ConfirmFriendRequestReq) error {
	if err := req.Validate(); err != nil {
		return responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error())
	}

	return s.confirmFriendRequest(ctx, req)
}

func (s *FriendService) RejectFriend(ctx context.Context, req *friendpb.ConfirmFriendRequestReq) error {
	if err := req.Validate(); err != nil {
		return responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error())
	}

	return s.confirmFriendRequest(ctx, req)
}

func (s *FriendService) confirmFriendRequest(ctx context.Context, req *friendpb.ConfirmFriendRequestReq) error {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).ConfirmFriendRequest(ctx, req)
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}

/*
 * Friend Logic
 */

func (s *FriendService) ListUserRelation(ctx context.Context, req *friendpb.QueryFriendListRequest) (
	[]*friendpb.Friend, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).QueryFriendList(ctx, req)
	if err != nil {
		return nil, err
	}

	if rsp.Response.Success() {
		return rsp.GetFriendList(), nil
	}

	return nil, rsp.GetResponse()
}

func (s *FriendService) BlockFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_BLOCKED)
}

func (s *FriendService) UnblockFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_UNBLOCKED)
}

func (s *FriendService) DeleteFriend(ctx context.Context, req *friendpb.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_STRANGER)
}

func (s *FriendService) updateFriendStatus(ctx context.Context, req *friendpb.BaseFriendRequest, status friendpb.FriendStatus) error { // nolint: lll
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return err
	}

	updateReq := &friendpb.UpdateFriendStatusRequest{
		Info:   req,
		Status: status,
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).UpdateFriendStatus(ctx, updateReq)
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}
