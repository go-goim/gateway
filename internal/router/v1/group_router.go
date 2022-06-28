package v1

import (
	"github.com/gin-gonic/gin"

	responsepb "github.com/go-goim/api/transport/response"
	grouppb "github.com/go-goim/api/user/group/v1"
	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/request"
	"github.com/go-goim/core/pkg/response"
	"github.com/go-goim/core/pkg/router"
	"github.com/go-goim/gateway/internal/service"
)

type GroupRouter struct {
	router.Router
}

func NewGroupRouter() *GroupRouter {
	return &GroupRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *GroupRouter) Load(router *gin.RouterGroup) {
	router.Use(mid.AuthJwtCookie)
	router.GET("/get", r.getGroup)
	router.GET("/list", r.listGroup)
	router.POST("/create", r.createGroup)
	router.POST("/update", r.updateGroup)
	router.POST("/delete", r.deleteGroup)
	router.POST("/join", r.joinGroup)
	router.POST("/leave", r.leaveGroup)
}

// @Summary 获取群组信息
// @Description 获取群组信息
// @Tags [gateway]群组
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Authorization header string true "token"
// @Param g_id query string true "群组ID"
// @Param with_member query bool false "是否获取群组成员"
// @Success 200 {object} grouppb.Group "Success"
// @Failure 400 {object} responsepb.BaseResponse "Bad Request"
// @Router /gateway/v1/group/get [get]
func (r *GroupRouter) getGroup(c *gin.Context) {
	var req struct {
		GroupID     string `form:"gid" binding:"required"`
		WithMembers bool   `form:"with_members"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	group, err := service.GetGroupService().GetGroup(mid.GetContext(c), &grouppb.GetGroupRequest{
		Gid:         req.GroupID,
		WithMembers: req.WithMembers,
	})
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, group)
}

// @Summary 获取群组列表
// @Description 获取群组列表
// @Tags [gateway]群组
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Authorization header string true "token"
// @Param uid query string true "用户ID"
// @Param page query int32 false "页码"
// @Param size query int32 false "每页数量"
// @Success 200 {object} response.Response{data=[]grouppb.Group} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /gateway/v1/group/list [get]
func (r *GroupRouter) listGroup(c *gin.Context) {
	var req struct {
		Uid  string `form:"uid" binding:"required"`
		Page int32  `form:"page"`
		Size int32  `form:"size"` // todo support default size
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	list, err := service.GetGroupService().ListGroup(mid.GetContext(c), &grouppb.ListGroupsRequest{
		Uid:      req.Uid,
		Page:     req.Page,
		PageSize: req.Size,
	})
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, list, response.SetTotal(len(list)))
}

// @Summary 创建群组
// @Description 创建群组
// @Tags [gateway]群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.CreateGroupRequest true "创建群组请求"
// @Success 200 {object} grouppb.Group "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /gateway/v1/group/create [post]
func (r *GroupRouter) createGroup(c *gin.Context) {
	var req = &grouppb.CreateGroupRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	group, err := service.GetGroupService().CreateGroup(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, group)
}

// @Summary 更新群组
// @Description 更新群组
// @Tags [gateway]群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.UpdateGroupRequest true "更新群组请求"
// @Success 200 {object} grouppb.Group "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /gateway/v1/group/update [post]
func (r *GroupRouter) updateGroup(c *gin.Context) {
	var req = &grouppb.UpdateGroupRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	group, err := service.GetGroupService().UpdateGroup(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, group)
}

// @Summary 删除群组
// @Description 删除群组
// @Tags [gateway]群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.DeleteGroupRequest true "删除群组请求"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /gateway/v1/group/delete [post]
func (r *GroupRouter) deleteGroup(c *gin.Context) {
	var req = &grouppb.DeleteGroupRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	err := service.GetGroupService().DeleteGroup(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, nil)
}

// @Summary 加入群组
// @Description 加入群组
// @Tags [gateway]群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.AddGroupMemberRequest true "加入群组请求"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /gateway/v1/group/join [post]
func (r *GroupRouter) joinGroup(c *gin.Context) {
	var req = &grouppb.AddGroupMemberRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	cnt, err := service.GetGroupService().AddGroupMember(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{"count": cnt})
}

// @Summary 退出群组
// @Description 退出群组
// @Tags [gateway]群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.RemoveGroupMemberRequest true "退出群组请求"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /gateway/v1/group/leave [post]
func (r *GroupRouter) leaveGroup(c *gin.Context) {
	var req = &grouppb.RemoveGroupMemberRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ErrorResp(c, responsepb.Code_InvalidParams.BaseResponseWithError(err))
		return
	}

	cnt, err := service.GetGroupService().RemoveGroupMember(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{"count": cnt})
}
