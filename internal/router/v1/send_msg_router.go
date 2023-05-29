package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/go-goim/api/errors"
	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/router"
	"github.com/go-goim/core/pkg/web/response"
	"github.com/go-goim/gateway/internal/dto"

	"github.com/go-goim/gateway/internal/service"
)

type MsgRouter struct {
	router.Router
}

func NewMsgRouter() *MsgRouter {
	return &MsgRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *MsgRouter) Load(g *gin.RouterGroup) {
	g.Use(mid.AuthJwt)
	offline := NewOfflineMessageRouter()
	offline.Load(g.Group("/offline"))

	g.POST("/send_msg", r.handleSendSingleUserMsg)
	g.POST("/broadcast", r.handleSendBroadcastMsg)
}

// @Summary 发送单聊消息
// @Description 发送单聊消息
// @Tags message
// @Accept  json
// @Produce  json
// @Param   Authorization header string true "token"
// @Param   req body dto.SendMessageReq true "req"
// @Success 200 {object} response.Response{data=dto.SendMessageResp} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /message/send_msg [post]
func (r *MsgRouter) handleSendSingleUserMsg(c *gin.Context) {
	req := new(dto.SendMessageReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, errors.ErrorCode_InvalidParams.WithError(err))
		return
	}

	rsp, err := service.GetSendMessageService().SendMessage(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}

// @Summary 发送广播消息
// @Description 发送广播消息
// @Tags message
// @Accept  json
// @Produce  json
// @Param   Authorization header string true "token"
// @Param   req body dto.SendMessageReq true "req"
// @Success 200 {object} response.Response{data=dto.SendMessageResp} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /message/broadcast [post]
func (r *MsgRouter) handleSendBroadcastMsg(c *gin.Context) {
	req := new(dto.SendMessageReq)
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, errors.ErrorCode_InvalidParams.WithError(err))
		return
	}

	rsp, err := service.GetSendMessageService().Broadcast(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, rsp)
}
