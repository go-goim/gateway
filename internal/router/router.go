package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/router"

	v1 "github.com/go-goim/gateway/internal/router/v1"
)

type rootRouter struct {
	router.Router
}

func newRootRouter() *rootRouter {
	r := &rootRouter{
		Router: &router.BaseRouter{},
	}

	r.init()
	return r
}

func (r *rootRouter) init() {
	r.Register("/v1", v1.NewRouter())
}

func RegisterRouter(g *gin.RouterGroup) {
	g.GET("/ping", ping)
	g.Use(mid.PagingHandler)
	r := newRootRouter()
	r.Load(g)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
