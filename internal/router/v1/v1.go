package v1

import (
	"github.com/go-goim/core/pkg/router"
)

type Router struct {
	router.Router
}

func NewRouter() *Router {
	r := &Router{
		Router: &router.BaseRouter{},
	}
	r.Init()
	return r
}

func (r *Router) Init() {
	r.Register("/user", NewUserRouter())
	r.Register("/group", NewGroupRouter())
	r.Register("/message", NewMsgRouter())
	r.Register("/discovery", NewDiscoverRouter())
}
