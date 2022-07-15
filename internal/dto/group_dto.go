package dto

import (
	grouppb "github.com/go-goim/api/user/group/v1"
	"github.com/go-goim/core/pkg/model"
)

type GetGroupRequest struct {
	UID         *model.ID `form:"-"`
	GID         *model.ID `form:"gid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	WithMembers bool      `form:"with_members" example:"true"`
	// WithInfo valid only when withMembers is true
	WithInfo bool `form:"with_info" example:"true"`
}

func (r *GetGroupRequest) ToPb() *grouppb.GetGroupRequest {
	pb := &grouppb.GetGroupRequest{}
	pb.Gid = r.GID.Int64()
	pb.WithMembers = r.WithMembers
	pb.WithInfo = r.WithInfo

	return pb
}

type CreateGroupRequest struct {
	UID     *model.ID   `json:"-"`
	Name    string      `json:"name" validate:"required,max=32" example:"test"`
	Desc    string      `json:"desc" validate:"omitempty,max=128" example:"test"`
	Avatar  string      `json:"avatar" validate:"omitempty,url" example:"https://example.com/avatar.png"`
	Members []*model.ID `json:"members" validate:"required,min=2,max=20" swaggertype:"array,string" example:"av8FMdRdcb,av8FMdRdcc"` //nolint:lll
}

func (r *CreateGroupRequest) ToPb() *grouppb.CreateGroupRequest {
	pb := &grouppb.CreateGroupRequest{}
	pb.Name = r.Name
	pb.Description = r.Desc
	pb.Avatar = r.Avatar
	pb.MembersUid = make([]int64, len(r.Members))
	for i, m := range r.Members {
		pb.MembersUid[i] = m.Int64()
	}

	return pb
}

type UpdateGroupRequest struct {
	UID    *model.ID `json:"-"`
	GID    *model.ID `json:"gid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	Name   *string   `json:"name" validate:"omitempty,max=32" example:"test"`
	Desc   *string   `json:"desc" validate:"omitempty,max=128" example:"test"`
	Avatar *string   `json:"avatar" validate:"omitempty,url" example:"https://www.example.com/avatar.png"`
}

func (r *UpdateGroupRequest) ToPb() *grouppb.UpdateGroupRequest {
	pb := &grouppb.UpdateGroupRequest{}
	pb.OwnerUid = r.UID.Int64()
	pb.Gid = r.GID.Int64()
	pb.Name = r.Name
	pb.Description = r.Desc
	pb.Avatar = r.Avatar

	return pb
}

type DeleteGroupRequest struct {
	UID *model.ID `json:"-"`
	GID *model.ID `json:"gid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
}

func (r *DeleteGroupRequest) ToPb() *grouppb.DeleteGroupRequest {
	pb := &grouppb.DeleteGroupRequest{}
	pb.OwnerUid = r.UID.Int64()
	pb.Gid = r.GID.Int64()

	return pb
}

type ChangeGroupMemberRequest struct {
	UID  *model.ID   `json:"-"`
	GID  *model.ID   `json:"gid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	UIDs []*model.ID `json:"uids" validate:"required,min=1,max=20" swaggertype:"array,string" example:"av8FMdRdcb,av8FMdRdcc"` //nolint:lll
}

func (r *ChangeGroupMemberRequest) ToPb() *grouppb.ChangeGroupMemberRequest {
	pb := &grouppb.ChangeGroupMemberRequest{}
	pb.OwnerUid = r.UID.Int64()
	pb.Gid = r.GID.Int64()
	pb.Uids = make([]int64, len(r.UIDs))
	for i, u := range r.UIDs {
		pb.Uids[i] = u.Int64()
	}

	return pb
}

type ChangeGroupMemberResponse struct {
	Count int `json:"count" example:"1"`
}

func ChangeGroupMemberResponseFromPb(pb *grouppb.ChangeGroupMemberResponse) *ChangeGroupMemberResponse {
	return &ChangeGroupMemberResponse{
		Count: int(pb.Count),
	}
}

type GroupMember struct {
	GID  *model.ID `json:"gid" swaggertype:"string" example:"av8FMdRdcb"`
	UID  *model.ID `json:"uid" swaggertype:"string" example:"av8FMdRdcb"`
	User *User     `json:"user,omitempty"` // only when withMembers is true and withInfo is true
	// 0: normal, 1: silent
	Status int32 `json:"status" example:"1"`
	// 0: owner, 1: member
	Type int32 `json:"type" example:"1"`
}

func GroupMemberFromPb(pb *grouppb.GroupMember) *GroupMember {
	return &GroupMember{
		GID:    model.NewID(pb.Gid),
		UID:    model.NewID(pb.Uid),
		User:   UserFromPb(pb.User),
		Status: int32(pb.Status),
		Type:   int32(pb.Type),
	}
}

type Group struct {
	GID         *model.ID      `json:"gid" swaggertype:"string" example:"av8FMdRdcb"`
	Name        string         `json:"name" example:"test"`
	Desc        string         `json:"desc" example:"test"`
	Avatar      string         `json:"avatar" example:"https://example.com/avatar.png"`
	OwnerUID    *model.ID      `json:"owner_uid" swaggertype:"string" example:"av8FMdRdcb"`
	Owner       *GroupMember   `json:"owner,omitempty"`
	Members     []*GroupMember `json:"members,omitempty"`
	MaxMembers  int32          `json:"max_member" example:"20"`
	MemberCount int32          `json:"member_count" example:"2"`
	Status      int32          `json:"status" example:"1"`
}

func GroupFromPb(pb *grouppb.Group) *Group {
	return &Group{
		GID:         model.NewID(pb.Gid),
		Name:        pb.Name,
		Desc:        pb.Description,
		Avatar:      pb.Avatar,
		OwnerUID:    model.NewID(pb.OwnerUid),
		Owner:       GroupMemberFromPb(pb.Owner),
		Members:     make([]*GroupMember, len(pb.Members)),
		MaxMembers:  pb.MaxMembers,
		MemberCount: pb.MemberCount,
		Status:      int32(pb.Status),
	}
}

func GroupsFromPb(pb []*grouppb.Group) []*Group {
	groups := make([]*Group, len(pb))
	for i, g := range pb {
		groups[i] = GroupFromPb(g)
	}
	return groups
}
