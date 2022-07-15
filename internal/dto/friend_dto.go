package dto

import (
	friendv1 "github.com/go-goim/api/user/friend/v1"
	sessionv1 "github.com/go-goim/api/user/session/v1"
	"github.com/go-goim/core/pkg/model"
)

type AddFriendResult struct {
	FriendRequest *FriendRequest `json:"friendRequest"`
	Status        int32          `json:"status"`
}

func AddFriendResultFromPb(pb *friendv1.AddFriendResult) *AddFriendResult {
	return &AddFriendResult{
		FriendRequest: FriendRequestFromPb(pb.FriendRequest),
		Status:        int32(pb.Status),
	}
}

type BaseFriendRequest struct {
	UID       *model.ID `json:"-"` // no validation
	FriendUID *model.ID `json:"friendUid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
}

func (r *BaseFriendRequest) ToPb() *friendv1.BaseFriendRequest {
	pb := &friendv1.BaseFriendRequest{}
	pb.Uid = r.UID.Int64()
	pb.FriendUid = r.FriendUID.Int64()

	return pb
}

type ConfirmFriendRequestRequest struct {
	UID             *model.ID `json:"-"`
	FriendRequestID int64     `json:"friendRequestId" validate:"required" example:"2"`
	Action          int32     `json:"-"`
}

func (r *ConfirmFriendRequestRequest) ToPb() *friendv1.ConfirmFriendRequestRequest {
	pb := &friendv1.ConfirmFriendRequestRequest{}
	pb.Uid = r.UID.Int64()
	pb.FriendRequestId = r.FriendRequestID
	pb.Action = friendv1.ConfirmFriendRequestAction(r.Action)

	return pb
}

type QueryFriendRequestListRequest struct {
	UID    *model.ID `json:"-" form:"-"`
	Status int32     `form:"status" validate:"required,oneof=0 1" example:"0"`
}

func (r *QueryFriendRequestListRequest) ToPb() *friendv1.QueryFriendRequestListRequest {
	pb := &friendv1.QueryFriendRequestListRequest{}
	pb.Uid = r.UID.Int64()
	pb.Status = friendv1.FriendRequestStatus(r.Status)

	return pb
}

type UpdateFriendStatusRequest struct {
	UID       *model.ID `json:"uid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	FriendUID *model.ID `json:"friendUid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	Status    int32     `json:"status" validate:"required,oneof=0 1 2 3" example:"0"`
}

func (r *UpdateFriendStatusRequest) ToPb() *friendv1.UpdateFriendStatusRequest {
	pb := &friendv1.UpdateFriendStatusRequest{}
	pb.Info = &friendv1.BaseFriendRequest{
		Uid:       r.UID.Int64(),
		FriendUid: r.FriendUID.Int64(),
	}
	pb.Status = friendv1.FriendStatus(r.Status)

	return pb
}

type CheckSendMessageAbilityRequest struct {
	FromUID     *model.ID `json:"fromUid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	ToUID       *model.ID `json:"toUid" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	SessionType int32     `json:"sessionType" validate:"required,gte=0,lte=255" example:"0"`
}

func (r *CheckSendMessageAbilityRequest) ToPb() *friendv1.CheckSendMessageAbilityRequest {
	pb := &friendv1.CheckSendMessageAbilityRequest{}
	pb.FromUid = r.FromUID.Int64()
	pb.ToUid = r.ToUID.Int64()
	pb.SessionType = sessionv1.SessionType(r.SessionType)

	return pb
}

type FriendRequest struct {
	ID           int64     `json:"id" example:"1"`
	UID          *model.ID `json:"uid" swaggertype:"string" example:"av8FMdRdcb"`
	FriendUID    *model.ID `json:"friendUid" swaggertype:"string" example:"av8FMdRdcb"`
	FriendName   string    `json:"friendName" example:"friendName"`
	FriendAvatar string    `json:"friendAvatar" example:"https://www.example.com/friendAvatar.png"`
	// 0: pending, 1: accepted, 2: rejected
	Status    int32 `json:"status" example:"0"`
	CreatedAt int64 `json:"createdAt" example:"1579098983"`
	UpdatedAt int64 `json:"updatedAt" example:"1579098983"`
}

func FriendRequestFromPb(pb *friendv1.FriendRequest) *FriendRequest {
	return &FriendRequest{
		ID:           pb.Id,
		UID:          model.NewID(pb.Uid),
		FriendUID:    model.NewID(pb.FriendUid),
		FriendName:   pb.FriendName,
		FriendAvatar: pb.FriendAvatar,
		Status:       int32(pb.Status),
		CreatedAt:    pb.CreatedAt,
		UpdatedAt:    pb.UpdatedAt,
	}
}

func FriendRequestListFromPb(pb []*friendv1.FriendRequest) []*FriendRequest {
	var list []*FriendRequest
	for _, v := range pb {
		list = append(list, FriendRequestFromPb(v))
	}
	return list
}

type Friend struct {
	UID          *model.ID `json:"uid" swaggertype:"string" example:"av8FMdRdcb"`
	FriendUID    *model.ID `json:"friendUid" swaggertype:"string" example:"av8FMdRdcb"`
	FriendName   string    `json:"friendName" example:"friendName"`
	FriendAvatar string    `json:"friendAvatar" example:"https://www.example.com/friendAvatar.png"`
	// 0: friend, 1: stranger, 2: blacklist
	Status    int32 `json:"status" example:"0"`
	CreatedAt int64 `json:"createdAt" example:"1579098983"`
	UpdatedAt int64 `json:"updatedAt" example:"1579098983"`
}

func FriendFromPb(pb *friendv1.Friend) *Friend {
	return &Friend{
		UID:          model.NewID(pb.Uid),
		FriendUID:    model.NewID(pb.FriendUid),
		FriendName:   pb.FriendName,
		FriendAvatar: pb.FriendAvatar,
		Status:       int32(pb.Status),
		CreatedAt:    pb.CreatedAt,
		UpdatedAt:    pb.UpdatedAt,
	}
}
