package v1

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/response"
	"github.com/go-goim/core/pkg/router"

	"github.com/go-goim/gateway/internal/service"
)

type DiscoverRouter struct {
	router.Router
}

func NewDiscoverRouter() *DiscoverRouter {
	return &DiscoverRouter{
		Router: &router.BaseRouter{},
	}
}

func (r *DiscoverRouter) Load(g *gin.RouterGroup) {
	g.GET("/discover", mid.AuthJwtCookie, r.handleDiscoverPushServer)
}

// @Summary 获取推送服务器
// @Description 获取推送服务器 IP
// @Tags [gateway]discover
// @Produce  json
// @Param   token query string true "token"
// @Success 200 {string} string "success"
// @Failure 200 {object} response.Response "errCode"
// @Failure 401 {null} null "unauthorized"
// @Router /gateway/v1/discovery/discover [get]
func (r *DiscoverRouter) handleDiscoverPushServer(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		response.ErrorResp(c, fmt.Errorf("uid is empty"))
		return
	}

	serverIP, err := service.LoadMatchedPushServer(context.Background())
	if err != nil {
		response.ErrorResp(c, err)
		return
	}

	response.SuccessResp(c, gin.H{
		"server_ip": serverIP,
	})
}
