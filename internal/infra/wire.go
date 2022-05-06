package infra

import (
	"github.com/google/wire"
	"github.com/yuridevx/app-example/internal/domain"
)

var Set = wire.NewSet(
	NewApp,
	NewRelic,
	NewLogger,
	NewPyroscope,
	NewHealth,
	wire.Bind(new(domain.HealthService), new(*Health)),
	NewClock,
	wire.Bind(new(domain.Clock), new(*Clock)),
)
