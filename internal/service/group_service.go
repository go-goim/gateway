package service

import (
	"context"

	grouppb "github.com/go-goim/api/user/group/v1"
	"github.com/go-goim/core/pkg/types"
	"github.com/go-goim/core/pkg/web"
	"github.com/go-goim/gateway/internal/dto"
)

type GroupService struct {
}

var (
	groupService = &GroupService{}
)

func GetGroupService() *GroupService {
	return groupService
}

func (s *GroupService) GetGroup(ctx context.Context, req *dto.GetGroupRequest) (*dto.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).GetGroup(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return dto.GroupFromPb(rsp.GetGroup()), nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, req *dto.UpdateGroupRequest) (*dto.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).UpdateGroup(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return dto.GroupFromPb(rsp.GetGroup()), nil
}

func (s *GroupService) CreateGroup(ctx context.Context, req *dto.CreateGroupRequest) (*dto.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).CreateGroup(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return dto.GroupFromPb(rsp.GetGroup()), nil
}

func (s *GroupService) ListGroup(ctx context.Context, uid *types.ID, paging *web.Paging) ([]*dto.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).ListGroups(ctx, &grouppb.ListGroupsRequest{
		Uid:      uid.Int64(),
		Page:     int32(paging.Page),
		PageSize: int32(paging.PageSize),
	})
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return dto.GroupsFromPb(rsp.GetGroups()), nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, req *dto.DeleteGroupRequest) error {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).DeleteGroup(ctx, req.ToPb())
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}

func (s *GroupService) AddGroupMember(ctx context.Context, req *dto.ChangeGroupMemberRequest) (int, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return 0, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).AddGroupMember(ctx, req.ToPb())
	if err != nil {
		return 0, err
	}

	if !rsp.Response.Success() {
		return 0, rsp.Response
	}

	return int(rsp.GetCount()), nil
}

func (s *GroupService) RemoveGroupMember(ctx context.Context, req *dto.ChangeGroupMemberRequest) (int, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return 0, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).RemoveGroupMember(ctx, req.ToPb())
	if err != nil {
		return 0, err
	}

	if !rsp.Response.Success() {
		return 0, rsp.Response
	}

	return int(rsp.GetCount()), nil
}
