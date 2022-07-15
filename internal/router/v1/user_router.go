package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/go-goim/gateway/internal/dto"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/router"
	"github.com/go-goim/core/pkg/web/response"

	"github.com/go-goim/gateway/internal/service"
)

type UserRouter struct {
	router.Router
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *UserRouter) Load(router *gin.RouterGroup) {
	friend := NewFriendRouter()
	friend.Load(router.Group("/friend", mid.AuthJwt))

	auth := router.Group("", mid.AuthJwt)
	{
		auth.GET("/query", r.queryUser)
		auth.POST("/update", r.updateUserInfo)
	}

	// no auth
	router.POST("/login", r.login)
	router.POST("/register", r.register)
}

// @Summary 查询用户信息
// @Description 查询用户信息
// @Tags 用户
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Authorization header string true "token"
// @Param email query string false "email"
// @Param Phone query string false "phone"
// @Success 200 {object} response.Response{data=dto.User}
// @Failure 400 {object} response.Response
// @Router /user/query [get]
func (r *UserRouter) queryUser(c *gin.Context) {
	var req = &dto.QueryUserRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().QueryUserInfo(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, user)
}

// @Summary 登录
// @Description 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param   req body dto.UserLoginRequest true "req"
// @Success 200 {object} response.Response{data=dto.User}
// @Header  400 {string} Authorization "Bearer "
// @Router /user/login [post]
func (r *UserRouter) login(c *gin.Context) {
	var req = &dto.UserLoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().Login(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	if err = mid.SetJwtToHeader(c, user.UID.Int64()); err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, user)
}

// @Summary 注册
// @Description 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param   req body dto.CreateUserRequest true "req"
// @Success 200 {object} response.Response{data=dto.User}
// @Failure 400 {object} response.Response
// @Router /user/register [post]
func (r *UserRouter) register(c *gin.Context) {
	var req = &dto.CreateUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().Register(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, user)
}

// @Summary 更新用户信息
// @Description 更新用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param   req body dto.UpdateUserRequest true "req"
// @Success 200 {object} response.Response{data=dto.User}
// @Failure 400 {object} response.Response
// @Router /user/update [post]
func (r *UserRouter) updateUserInfo(c *gin.Context) {
	var req = &dto.UpdateUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ErrorResp(c, err)
		return
	}

	user, err := service.GetUserService().UpdateUser(mid.GetContext(c), req)
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, user)
}
