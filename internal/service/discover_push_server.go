package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-goim/core/pkg/registry"

	"github.com/go-goim/gateway/internal/app"
)

func LoadMatchedPushServer(ctx context.Context) (string, error) {
	list, err := registry.GetService(ctx, app.GetApplication().Config.SrvConfig.PushService)
	if err != nil {
		return "", err
	}

	if len(list) == 0 {
		return "", fmt.Errorf("push service not found")
	}

	var eps []string
	for _, instance := range list {
		for _, ep := range instance.Endpoints {
			if strings.HasPrefix(ep, "http") {
				eps = append(eps, ep)
			}
		}
	}
	if len(eps) == 0 {
		return "", fmt.Errorf("no matched service")
	}

	return eps[rand.Int()%len(eps)], nil // nolint:gosec
}
