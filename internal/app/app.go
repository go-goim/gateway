package app

import (
	"github.com/go-goim/core/pkg/app"
	"github.com/go-goim/core/pkg/cmd"
	"github.com/go-goim/core/pkg/registry"
	"github.com/go-goim/core/pkg/util/snowflake"
)

type Application struct {
	*app.Application
	// add own fields here
}

var (
	application *Application
	nodeBit     int64
)

func init() {
	// TODO: use hostname and ip to get nodeBit from registry
	cmd.GlobalFlagSet.Int64VarP(&nodeBit, "node-bit", "b", 1, "node bit")
}

func InitApplication() (*Application, error) {
	snowflake.SetDefaultNode(nodeBit)
	// do some own biz logic if needed
	a, err := app.InitApplication()
	if err != nil {
		return nil, err
	}
	application = &Application{Application: a}

	return application, nil
}

func GetRegister() registry.RegisterDiscover {
	return application.Register
}

func GetApplication() *Application {
	app.AssertApplication()
	return application
}
