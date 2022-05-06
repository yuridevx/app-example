package logic

import (
	"github.com/google/wire"
	v1 "github.com/yuridevx/app-example/api/v1/schedule"
	"github.com/yuridevx/app-example/internal/domain"
	"github.com/yuridevx/app-example/internal/logic/runner"
	"github.com/yuridevx/app-example/internal/logic/schedule"
	"github.com/yuridevx/app-example/internal/logic/server"
)

var Set = wire.NewSet(
	schedule.NewController,
	wire.Bind(new(v1.ScheduleServer), new(*schedule.Controller)),
	server.NewGrpc,
	server.NewGateway,
	runner.NewRunner,
	wire.Bind(new(domain.Runner), new(*runner.Runner)),
)
