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
	router.Use(mid.AuthJwt)
	router.GET("/get", r.getGroup)
	router.GET("/list", r.listGroup)
	router.POST("/create", r.createGroup)
	router.POST("/update", r.updateGroup)
	router.POST("/delete", r.deleteGroup)
	router.POST("/member/add", r.addMember)
	router.POST("/member/remove", r.removeMember) // kick by manager
	router.POST("/member/leave", r.leaveGroup)    // leave group by themselves
}

// @Summary 获取群组信息
// @Description 获取群组信息
// @Tags 群组
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Authorization header string true "token"
// @Param gid query string true "群组ID"
// @Param with_members query bool false "是否获取群组成员"
// @Success 200 {object} grouppb.Group "Success"
// @Failure 400 {object} responsepb.BaseResponse "Bad Request"
// @Router /group/get [get]
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
// @Tags 群组
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Authorization header string true "token"
// @Param page query int32 false "页码"
// @Param size query int32 false "每页数量"
// @Success 200 {object} response.Response{data=[]grouppb.Group} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /group/list [get]
func (r *GroupRouter) listGroup(c *gin.Context) {
	var req struct {
		Page int32 `form:"page"`
		Size int32 `form:"size"` // todo support default size
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	list, err := service.GetGroupService().ListGroup(mid.GetContext(c), &grouppb.ListGroupsRequest{
		Uid:      mid.GetUID(c),
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
// @Tags 群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.CreateGroupRequest true "创建群组请求"
// @Success 200 {object} grouppb.Group "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /group/create [post]
func (r *GroupRouter) createGroup(c *gin.Context) {
	var req = &grouppb.CreateGroupRequest{}
	if err := c.ShouldBindWith(req, &request.PbJSONBinding{}); err != nil {
		response.ErrorResp(c, err)
		return
	}

	req.OwnerUid = mid.GetUID(c)
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
// @Tags 群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.UpdateGroupRequest true "更新群组请求"
// @Success 200 {object} grouppb.Group "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /group/update [post]
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
// @Tags 群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.DeleteGroupRequest true "删除群组请求"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /group/delete [post]
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

// @Summary 添加群组成员
// @Description 任何群成员都可以添加群组成员
// @Tags 群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.AddGroupMemberRequest true "加入群组请求"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /group/member/invite [post]
func (r *GroupRouter) addMember(c *gin.Context) {
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

// @Summary 删除群组成员
// @Description 群管理员可以删除群组成员
// @Tags 群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body grouppb.RemoveGroupMemberRequest true "删除群组成员请求"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /group/member/remove [post]
func (r *GroupRouter) removeMember(c *gin.Context) {
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

type leaveGroupRequest struct {
	GroupId string `json:"gid" validate:"required"`
}

// @Summary 退出群组
// @Description 群组成员可以退出群组
// @Tags 群组
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param body body leaveGroupRequest true "退出群组请求"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Router /group/member/leave [post]
func (r *GroupRouter) leaveGroup(c *gin.Context) {
	var req = &leaveGroupRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	cnt, err := service.GetGroupService().RemoveGroupMember(mid.GetContext(c), &grouppb.RemoveGroupMemberRequest{
		Gid: req.GroupId,
		Uid: []string{mid.GetUID(c)},
	})
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{"count": cnt})
}
