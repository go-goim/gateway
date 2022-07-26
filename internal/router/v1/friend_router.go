package v1

import (
	"github.com/gin-gonic/gin"

	responsepb "github.com/go-goim/api/transport/response"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/gateway/internal/dto"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/router"
	"github.com/go-goim/core/pkg/web/response"

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
// @Tags 好友
// @Produce json
// @Param Authorization header string true "token"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {object} response.Response{data=[]dto.Friend} "Success"
// @Failure 400 {object} response.Response{} "err"
// @Router /user/friend/list [get]
func (r *FriendRouter) listRelation(c *gin.Context) {
	list, err := service.GetFriendService().ListUserRelation(mid.GetContext(c),
		mid.GetUID(c), mid.GetPaging(c))
	if err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	response.SuccessResp(c, list, response.SetTotal(int32(len(list))))
}

// @Summary 获取好友请求列表
// @Description 获取好友请求列表
// @Tags 好友
// @Produce json
// @Param Authorization header string true "token"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param status query int false "status"
// @Success 200 {object} response.Response{data=[]dto.FriendRequest} "Success"
// @Failure 400 {object} response.Response{} "err"
// @Router /user/friend/request/list [get]
func (r *FriendRouter) listFriendRequest(c *gin.Context) {
	req := &dto.QueryFriendRequestListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	list, err := service.GetFriendService().QueryFriendRequestList(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{"list": list}, response.SetTotal(int32(len(list))))
}

// query user info before add friend
// @Summary 添加好友
// @Description 添加好友
// @Tags 好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body dto.BaseFriendRequest true "body"
// @Success 200 {object} response.Response{data=[]dto.FriendRequest} "Success"
// @Failure 400 {object} response.Response{} "error"
// @Router /user/friend/add [post]
func (r *FriendRouter) addFriend(c *gin.Context) {
	req := &dto.BaseFriendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	log.Info("add friend request", "uid", req.UID, "friend_uid", req.FriendUID)
	result, err := service.GetFriendService().AddFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, result)
}

// @Summary 删除好友
// @Description 删除好友
// @Tags 好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body dto.BaseFriendRequest true "body"
// @Success 200 {object} response.Response{} "Success"
// @Failure 400 {object} response.Response{} "error"
// @Router /user/friend/delete [post]
func (r *FriendRouter) deleteFriend(c *gin.Context) {
	req := &dto.BaseFriendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	err := service.GetFriendService().DeleteFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 接受好友请求
// @Description 接受好友请求
// @Tags 好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body dto.ConfirmFriendRequestRequest true "body"
// @Success 200 {object} response.Response{} "Success"
// @Failure 400 {object} response.Response{} "error"
// @Router /user/friend/accept [post]
func (r *FriendRouter) acceptFriend(c *gin.Context) {
	req := &dto.ConfirmFriendRequestRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	err := service.GetFriendService().AcceptFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 拒绝好友请求
// @Description 拒绝好友请求
// @Tags 好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body dto.ConfirmFriendRequestRequest true "body"
// @Success 200 {object} response.Response{} "Success"
// @Failure 400 {object} response.Response{} "error"
// @Router /user/friend/reject [post]
func (r *FriendRouter) rejectFriend(c *gin.Context) {
	req := &dto.ConfirmFriendRequestRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	err := service.GetFriendService().RejectFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 屏蔽好友
// @Description 屏蔽好友
// @Tags 好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body dto.BaseFriendRequest true "body"
// @Success 200 {object} response.Response{} "Success"
// @Failure 400 {object} response.Response{} "error"
// @Router /user/friend/block [post]
func (r *FriendRouter) blockFriend(c *gin.Context) {
	req := &dto.BaseFriendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	err := service.GetFriendService().BlockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}

// @Summary 取消屏蔽好友
// @Description 取消屏蔽好友
// @Tags 好友
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body dto.BaseFriendRequest true "body"
// @Success 200 {object} response.Response{} "Success"
// @Failure 400 {object} response.Response{} "error"
// @Router /user/friend/unblock [post]
func (r *FriendRouter) unblockFriend(c *gin.Context) {
	req := &dto.BaseFriendRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	req.UID = mid.GetUID(c)
	err := service.GetFriendService().UnblockFriend(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.OK(c)
}
