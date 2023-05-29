package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/go-goim/api/errors"
	"github.com/go-goim/gateway/internal/dto"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/router"
	"github.com/go-goim/core/pkg/web/response"

	"github.com/go-goim/gateway/internal/service"
)

type OfflineMessageRouter struct {
	router.Router
}

func NewOfflineMessageRouter() *OfflineMessageRouter {
	return &OfflineMessageRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *OfflineMessageRouter) Load(g *gin.RouterGroup) {
	g.GET("/query", r.handleQueryOfflineMessage)
}

// @Summary 查询离线消息
// @Description 查询离线消息
// @Tags offline_message
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Authorization header string true "token"
// @Param lastMessageId query integer true "lastMessageId"
// @Param onlyCount query boolean false "onlyCount"
// @Param page query integer false "page"
// @Param pageSize query integer false "pageSize"
// @Success 200 {object} response.Response{data=[]dto.Message} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /message/offline/query [get]
func (r *OfflineMessageRouter) handleQueryOfflineMessage(c *gin.Context) {
	req := new(dto.QueryOfflineMessageReq)
	if err := c.ShouldBindQuery(req); err != nil {
		response.ErrorResp(c, errors.ErrorCode_InvalidParams.WithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	req.Paging = mid.GetPaging(c)
	messages, cnt, err := service.GetOfflineMessageService().QueryOfflineMsg(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, messages, response.SetTotal(cnt), response.SetPaging(req.Paging))
}
