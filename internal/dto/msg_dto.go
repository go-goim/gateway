package dto

import (
	messagev1 "github.com/go-goim/api/message/v1"
	"github.com/go-goim/core/pkg/types"
	"github.com/go-goim/core/pkg/web"
)

type SendMessageReq struct {
	From        types.ID `json:"from" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	To          types.ID `json:"to" validate:"required" swaggertype:"string" example:"av8FMdRdcb"`
	SessionType int32    `json:"sessionType" validate:"required" example:"1"`
	SessionID   *string  `json:"sessionId" validate:"required" example:"1"`
	ContentType int32    `json:"contentType" validate:"required" example:"1"`
	Content     string   `json:"content" validate:"required" example:"hello"`
}

func (r *SendMessageReq) ToPb() *messagev1.SendMessageReq {
	pb := &messagev1.SendMessageReq{}
	pb.From = r.From.Int64()
	pb.To = r.To.Int64()
	pb.SessionType = messagev1.SessionType(r.SessionType)
	pb.SessionId = r.SessionID
	pb.ContentType = messagev1.MessageContentType(r.ContentType)
	pb.Content = r.Content

	return pb
}

type SendMessageResp struct {
	MessageID int64  `json:"messageId" example:"1"`
	SessionID string `json:"sessionId" example:"abc"`
}

type Message struct {
	MessageID   int64    `json:"messageId" example:"1"`
	From        types.ID `json:"from" swaggertype:"string" example:"av8FMdRdcb"`
	To          types.ID `json:"to" swaggertype:"string" example:"av8FMdRdcb"`
	SessionType int32    `json:"sessionType" example:"1"`
	SessionID   string   `json:"sessionId" example:"1"`
	ContentType int32    `json:"contentType" example:"1"`
	Content     string   `json:"content" example:"hello"`
	CreateTime  int64    `json:"createTime" example:"1579098983"`
}

func MessageFromPb(pb *messagev1.Message) *Message {
	return &Message{
		MessageID:   pb.GetMsgId(),
		From:        types.ID(pb.GetFrom()),
		To:          types.ID(pb.GetTo()),
		SessionType: int32(pb.GetSessionType()),
		SessionID:   pb.GetSessionId(),
		ContentType: int32(pb.GetContentType()),
		Content:     pb.GetContent(),
		CreateTime:  pb.GetCreateTime(),
	}
}

func MessagesFromPb(pb []*messagev1.Message) []*Message {
	msgs := make([]*Message, len(pb))
	for i, p := range pb {
		msgs[i] = MessageFromPb(p)
	}
	return msgs
}

type QueryOfflineMessageReq struct {
	UID           types.ID `form:"-"`
	LastMessageID int64    `form:"lastMessageId" validate:"required" example:"1"`
	OnlyCount     bool     `form:"onlyCount" example:"true"`
	*web.Paging
}

func (r *QueryOfflineMessageReq) ToPb() *messagev1.QueryOfflineMessageReq {
	pb := &messagev1.QueryOfflineMessageReq{}
	pb.Uid = r.UID.Int64()
	pb.LastMsgId = r.LastMessageID
	pb.OnlyCount = r.OnlyCount
	pb.Page = r.Page
	pb.PageSize = r.PageSize

	return pb
}
