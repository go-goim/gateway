package cmd

import (
	"context"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/go-goim/core/pkg/cmd"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/mid"

	// swagger
	_ "github.com/swaggo/swag"

	"github.com/go-goim/gateway/internal/app"
	"github.com/go-goim/gateway/internal/router"

	// register swagger
	_ "github.com/go-goim/gateway/docs"
)

var (
	jwtSecret = ""
)

func init() {
	cmd.GlobalFlagSet.StringVar(&jwtSecret, "jwt-secret", "", "jwt secret")
}

func Main() {
	if err := cmd.ParseFlags(); err != nil {
		panic(err)
	}

	if jwtSecret == "" {
		panic("jwt secret is empty")
	}
	mid.SetJwtHmacSecret(jwtSecret)

	application, err := app.InitApplication()
	if err != nil {
		log.Fatal("initApplication got err", "error", err)
	}

	g := gin.New()
	g.Use(gin.Recovery(), mid.Logger)
	// Is this necessary to set a prefix? Cause this service name is gateway already.
	router.RegisterRouter(g.Group("/gateway"))
	application.HTTPSrv.HandlePrefix("/", g)
	// register swagger
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err = application.Run(); err != nil {
		log.Error("application run got error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Error("graceful shutdown got error", "error", err)
	}
}
