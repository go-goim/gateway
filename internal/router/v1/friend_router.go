package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	responsepb "github.com/go-goim/api/transport/response"
	friendpb "github.com/go-goim/api/user/friend/v1"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/request"
	"github.com/go-goim/core/pkg/response"
	"github.com/go-goim/core/pkg/router"

	"github.com/go-goim/gateway/internal/service"
)

type FriendRouter struct {
	router.Router
}

func NewFriendRouter() *FriendRouter {
	return &FriendRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *FriendRouter) Load(g *gin.RouterGroup) {
	g.GET("/list", r.listRelation)
	g.GET("/request/list", r.listFriendRequest)

	g.POST("/add", r.addFriend)
	g.POST("/delete", r.deleteFriend)
	g.POST("/accept", r.acceptFriend)
	g.POST("/reject", r.rejectFriend)
	g.POST("/block", r.blockFriend)
	g.POST("/unblock", r.unblockFriend)
}

// @Summary 获取好友列表
// @Description 获取好友列表
// @Tags [gateway]好友
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {object} response.Response{data=[]friendpb.Friend} "Success"
// @Failure 400 {object} response.Response{} "err"
// @Router /gateway/v1/user/friend/list [get]
func (r *FriendRouter) listRelation(c *gin.Context) {
	// no need to check uid
	uid := c.GetString("uid")
	req := &friendpb.QueryFriendListRequest{
		Uid: uid,
	}

	list, err := service.GetFriendService().ListUserRelation(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, list, response.SetTotal(len(list)))
}

// @Summary 获取好友请求列表
// @Description 获取好友请求列表
// @Tags [gateway]好友
// @Produce json
// @Param Authorization header string true "token"
// @Success 200 {object} response.Response{} "Success"
// @Failure 400 {object} response.Response{} "err"
// @Router /gateway/v1/user/friend/request/list [get]
func (r *FriendRouter) listFriendRequest(c *gin.Context) {
	uid := c.GetString("uid")
	req := &friendpb.QueryFriendRequestListRequest{
		Uid: uid,
	}

	var status int
	statusStr := c.Query("status")
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	req.Status = friendpb.FriendRequestStatus(status)

	list, err := service.GetFriendService().QueryFriendRequestList(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{"list": list}, response.SetTotal(len(list)))
}

// query user info before add friend
// @Summary 添加好友
// @Description 添加好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.AddFriendRequest true "body"
// @Success 200 {object} friendpb.AddFriendResult
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/add [post]
func (r *FriendRouter) addFriend(c *gin.Context) {
	req := &friendpb.AddFriendRequest{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.Uid = mid.GetUID(c)
	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error()))
		return
	}

	result, err := service.GetFriendService().AddFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, result)
}

// @Summary 删除好友
// @Description 删除好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.BaseFriendRequest true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/delete [post]
func (r *FriendRouter) deleteFriend(c *gin.Context) {
	req := &friendpb.BaseFriendRequest{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetFriendService().DeleteFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 接受好友请求
// @Description 接受好友请求
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.ConfirmFriendRequestReq true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/accept [post]
func (r *FriendRouter) acceptFriend(c *gin.Context) {
	req := &friendpb.ConfirmFriendRequestReq{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.Uid = mid.GetUID(c)
	req.Action = friendpb.ConfirmFriendRequestAction_ACCEPT
	err := service.GetFriendService().AcceptFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 拒绝好友请求
// @Description 拒绝好友请求
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.ConfirmFriendRequestReq true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/reject [post]
func (r *FriendRouter) rejectFriend(c *gin.Context) {
	req := &friendpb.ConfirmFriendRequestReq{}
	if err := c.ShouldBindWith(req, request.NonValidatePbJSONBinding); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.Action = friendpb.ConfirmFriendRequestAction_REJECT
	req.Uid = mid.GetUID(c)
	err := service.GetFriendService().RejectFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 屏蔽好友
// @Description 屏蔽好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.BaseFriendRequest true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/block [post]
func (r *FriendRouter) blockFriend(c *gin.Context) {
	req := &friendpb.BaseFriendRequest{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetFriendService().BlockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 取消屏蔽好友
// @Description 取消屏蔽好友
// @Tags [gateway]好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body friendpb.BaseFriendRequest true "body"
// @Success 200 {object} responsepb.BaseResponse "Success"
// @Failure 400 {object} responsepb.BaseResponse "err"
// @Router /gateway/v1/user/friend/unblock [post]
func (r *FriendRouter) unblockFriend(c *gin.Context) {
	req := &friendpb.BaseFriendRequest{}
	if err := c.ShouldBindWith(req, request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	err := service.GetFriendService().UnblockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}
