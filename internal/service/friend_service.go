package service

import (
	"context"

	friendpb "github.com/go-goim/api/user/friend/v1"
	"github.com/go-goim/core/pkg/types"
	"github.com/go-goim/core/pkg/web"
	"github.com/go-goim/gateway/internal/dto"
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

func (s *FriendService) AddFriend(ctx context.Context, req *dto.BaseFriendRequest) (*dto.AddFriendResult, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).AddFriend(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return dto.AddFriendResultFromPb(rsp.GetResult()), nil
}

func (s *FriendService) QueryFriendRequestList(ctx context.Context, req *dto.QueryFriendRequestListRequest) (
	[]*dto.FriendRequest, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).QueryFriendRequestList(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return dto.FriendRequestListFromPb(rsp.GetFriendRequestList()), nil
}

func (s *FriendService) AcceptFriend(ctx context.Context, req *dto.ConfirmFriendRequestRequest) error {
	req.Action = int32(friendpb.ConfirmFriendRequestAction_ACCEPT)
	return s.confirmFriendRequest(ctx, req)
}

func (s *FriendService) RejectFriend(ctx context.Context, req *dto.ConfirmFriendRequestRequest) error {
	req.Action = int32(friendpb.ConfirmFriendRequestAction_REJECT)
	return s.confirmFriendRequest(ctx, req)
}

func (s *FriendService) confirmFriendRequest(ctx context.Context, req *dto.ConfirmFriendRequestRequest) error {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).ConfirmFriendRequest(ctx, req.ToPb())
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

func (s *FriendService) ListUserRelation(ctx context.Context, uid types.ID, paging *web.Paging) (
	[]*dto.Friend, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := friendpb.NewFriendServiceClient(cc).QueryFriendList(ctx, &friendpb.QueryFriendListRequest{
		Uid: uid.Int64(),
		// todo: paging
	})
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return dto.FriendsFromPb(rsp.GetFriendList()), nil
}

func (s *FriendService) BlockFriend(ctx context.Context, req *dto.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_BLOCKED)
}

func (s *FriendService) UnblockFriend(ctx context.Context, req *dto.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_UNBLOCKED)
}

func (s *FriendService) DeleteFriend(ctx context.Context, req *dto.BaseFriendRequest) error {
	return s.updateFriendStatus(ctx, req, friendpb.FriendStatus_STRANGER)
}

func (s *FriendService) updateFriendStatus(ctx context.Context, req *dto.BaseFriendRequest, status friendpb.FriendStatus) error { // nolint: lll
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return err
	}

	updateReq := &friendpb.UpdateFriendStatusRequest{
		Info:   req.ToPb(),
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
