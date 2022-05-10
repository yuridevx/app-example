//go:build wireinject
// +build wireinject

package wired

import (
	"github.com/google/wire"
	"github.com/yuridevx/app"
	"github.com/yuridevx/app-example/internal/config"
	"github.com/yuridevx/app-example/internal/domain"
	"github.com/yuridevx/app-example/internal/infra"
	"github.com/yuridevx/app-example/internal/logic"
	"github.com/yuridevx/app-example/internal/logic/server"
	"go.uber.org/zap"
)

type App struct {
	Builder app.Builder
	Logger  *zap.Logger
	Health  domain.HealthService
	// Commenting following field will prevent component from running
	Grpc    *server.Grpc
	Gateway *server.Gateway
}

func InitApp() (*App, error) {
	panic(
		wire.Build(
			config.Set,
			infra.Set,
			logic.Set,
			wire.Struct(new(App), "*"),
		),
	)
}
