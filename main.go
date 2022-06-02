//go:generate swag init -o ./docs --parseDependency --parseInternal --parseDepth 10
package main

import (
	"github.com/go-goim/gateway/cmd"
)

// @title GoIM.Gateway Swagger
// @version 1.0
// @description GoIM.Gateway 服务器接口文档
// @termsOfService http://go-goim.github.io/

// @contact.name Yusank
// @contact.url https://yusank.space
// @contact.email yusankurban@gmail.com

// @license.name MIT
// @license.url https://github.com/go-goim/core/blob/main/LICENSE

// @BasePath /gateway/v1
func main() {
	cmd.Main()
}
