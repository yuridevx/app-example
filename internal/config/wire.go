package config

import "github.com/google/wire"

var Set = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(*Config), "Relic", "Pyroscope"),
)
