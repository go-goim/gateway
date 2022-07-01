package service

import (
	"context"

	grouppb "github.com/go-goim/api/user/group/v1"
)

type GroupService struct {
}

var (
	groupService = &GroupService{}
)

func GetGroupService() *GroupService {
	return groupService
}

func (s *GroupService) GetGroup(ctx context.Context, req *grouppb.GetGroupRequest) (*grouppb.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).GetGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetGroup(), nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, req *grouppb.UpdateGroupRequest) (*grouppb.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).UpdateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetGroup(), nil
}

func (s *GroupService) CreateGroup(ctx context.Context, req *grouppb.CreateGroupRequest) (*grouppb.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetGroup(), nil
}

func (s *GroupService) ListGroup(ctx context.Context, req *grouppb.ListGroupsRequest) ([]*grouppb.Group, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).ListGroups(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.Response.Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetGroups(), nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, req *grouppb.DeleteGroupRequest) error {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).DeleteGroup(ctx, req)
	if err != nil {
		return err
	}

	if !rsp.Success() {
		return rsp
	}

	return nil
}

func (s *GroupService) AddGroupMember(ctx context.Context, req *grouppb.AddGroupMemberRequest) (int, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return 0, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).AddGroupMember(ctx, req)
	if err != nil {
		return 0, err
	}

	if !rsp.Response.Success() {
		return 0, rsp.Response
	}

	return int(rsp.GetAdded()), nil
}

func (s *GroupService) RemoveGroupMember(ctx context.Context, req *grouppb.RemoveGroupMemberRequest) (int, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return 0, err
	}

	rsp, err := grouppb.NewGroupServiceClient(cc).RemoveGroupMember(ctx, req)
	if err != nil {
		return 0, err
	}

	if !rsp.Response.Success() {
		return 0, rsp.Response
	}

	return int(rsp.GetRemoved()), nil
}
