package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"

	cgrpc "github.com/go-goim/core/pkg/conn/grpc"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/initialize"
	"github.com/go-goim/gateway/internal/app"
)

var (
	userServiceConnPool = &cgrpc.ConnPool{}
	msgServiceConnPool  = &cgrpc.ConnPool{}
)

func init() {
	initialize.Register(initialize.NewBasicInitializer("user_service", nil, func(_ context.Context) error {
		return initConnPool()
	}))
}

func initConnPool() error {
	var err error
	userServiceConnPool, err = cgrpc.NewConnPool(cgrpc.WithInsecure(),
		cgrpc.WithClientOption(
			grpc.WithEndpoint(fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.UserService)),
			grpc.WithDiscovery(app.GetApplication().Register),
			grpc.WithTimeout(time.Second*5),
			// grpc.WithOptions(ggrpc.WithBlock()),
		), cgrpc.WithPoolSize(2))
	if err != nil {
		return err
	}

	graceful.Register(func(_ context.Context) error {
		return userServiceConnPool.Release()
	})

	msgServiceConnPool, err = cgrpc.NewConnPool(cgrpc.WithInsecure(),
		cgrpc.WithClientOption(
			grpc.WithEndpoint(fmt.Sprintf("discovery://dc1/%s", app.GetApplication().Config.SrvConfig.MsgService)),
			grpc.WithDiscovery(app.GetApplication().Register),
			grpc.WithTimeout(time.Second*5),
			// grpc.WithOptions(ggrpc.WithBlock()),
		), cgrpc.WithPoolSize(2))
	if err != nil {
		return err
	}

	graceful.Register(func(_ context.Context) error {
		return msgServiceConnPool.Release()
	})

	return nil
}
