package dto

import (
	userv1 "github.com/go-goim/api/user/v1"
	"github.com/go-goim/core/pkg/types"
)

// Requests

type QueryUserRequest struct {
	// Email and Phone only one can be set
	Email *string `form:"email" validate:"omitempty,email" example:"user1@example.com"`
	Phone *string `form:"phone" validate:"omitempty,regexp=^1[3456789]\\d{9}$" example:"13800138000"`
}

func (r *QueryUserRequest) ToPb() *userv1.QueryUserRequest {
	pb := &userv1.QueryUserRequest{}
	if r.Email != nil {
		pb.Field = &userv1.QueryUserRequest_Email{Email: *r.Email}
	}
	if r.Phone != nil {
		pb.Field = &userv1.QueryUserRequest_Phone{Phone: *r.Phone}
	}

	return pb
}

type UserLoginRequest struct {
	// Email and Phone only one can be set
	Email     *string `json:"email" validate:"omitempty,email" example:"user1@example.com"`
	Phone     *string `json:"phone" validate:"omitempty,regexp=^1[3456789]\\d{9}$" example:"13800138000"`
	Password  string  `json:"password" validate:"required,min=6,max=20" example:"123456"`
	LoginType int32   `json:"loginType" validate:"required,oneof=0 1" example:"0"`
}

func (r *UserLoginRequest) ToPb() *userv1.UserLoginRequest {
	pb := &userv1.UserLoginRequest{}
	if r.Email != nil {
		pb.Field = &userv1.UserLoginRequest_Email{Email: *r.Email}
	}
	if r.Phone != nil {
		pb.Field = &userv1.UserLoginRequest_Phone{Phone: *r.Phone}
	}
	pb.Password = r.Password
	pb.LoginType = userv1.LoginType(r.LoginType)

	return pb
}

type CreateUserRequest struct {
	// Email and Phone only one can be set
	Email    *string `json:"email" validate:"omitempty,email" example:"user1@example.com"`
	Phone    *string `json:"phone" validate:"omitempty,regexp=^1[3456789]\\d{9}$" example:"13800138000"`
	Password string  `json:"password" validate:"required,min=6,max=20" example:"123456"`
	Name     string  `json:"name" validate:"required,min=2,max=32" example:"user1"`
}

func (r *CreateUserRequest) ToPb() *userv1.CreateUserRequest {
	pb := &userv1.CreateUserRequest{}
	if r.Email != nil {
		pb.Field = &userv1.CreateUserRequest_Email{Email: *r.Email}
	}
	if r.Phone != nil {
		pb.Field = &userv1.CreateUserRequest_Phone{Phone: *r.Phone}
	}
	pb.Password = r.Password
	pb.Name = r.Name

	return pb
}

type UpdateUserRequest struct {
	Email    *string `json:"email" validate:"omitempty,email" example:"user1@example.com"`
	Phone    *string `json:"phone" validate:"omitempty,regexp=^1[3456789]\\d{9}$" example:"13800138000"`
	Name     *string `json:"name" validate:"omitempty,min=2,max=32" example:"user1"`
	Password *string `json:"password" validate:"omitempty,min=6,max=20" example:"123456"`
	Avatar   *string `json:"avatar" validate:"omitempty,url" example:"https://www.example.com/avatar.png"`
}

func (r *UpdateUserRequest) ToPb() *userv1.UpdateUserRequest {
	pb := &userv1.UpdateUserRequest{}
	if r.Email != nil {
		pb.Email = *r.Email
	}
	if r.Phone != nil {
		pb.Phone = *r.Phone
	}
	if r.Name != nil {
		pb.Name = *r.Name
	}
	if r.Password != nil {
		pb.Password = *r.Password
	}
	if r.Avatar != nil {
		pb.Avatar = *r.Avatar
	}

	return pb
}

type User struct {
	UID         types.ID `json:"uid" swaggertype:"string" example:"av8FMdRdcb"`
	Name        string   `json:"name" example:"user1"`
	Avatar      string   `json:"avatar" example:"https://www.example.com/avatar.png"`
	Email       *string  `json:"email,omitempty" example:"abc@example.com"`
	Phone       *string  `json:"phone,omitempty" example:"13800138000"`
	ConnectURL  *string  `json:"connectUrl,omitempty" example:"ws://10.0.0.1:8080/ws"`
	LoginStatus int32    `json:"loginStatus" example:"0"`
}

func UserFromPb(pb *userv1.User) *User {
	if pb == nil {
		return nil
	}

	return &User{
		UID:         types.ID(pb.Uid),
		Name:        pb.Name,
		Avatar:      pb.Avatar,
		Email:       pb.Email,
		Phone:       pb.Phone,
		ConnectURL:  pb.PushServerIp,
		LoginStatus: int32(pb.LoginStatus),
	}
}
