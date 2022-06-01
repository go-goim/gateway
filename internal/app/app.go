package app

import (
	"github.com/go-goim/core/pkg/app"
	"github.com/go-goim/core/pkg/registry"
)

type Application struct {
	*app.Application
	// add own fields here
}

var (
	application *Application
)

func InitApplication() (*Application, error) {
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
