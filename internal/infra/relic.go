package infra

import (
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yuridevx/app-example/internal/config"
	"go.uber.org/zap"
)

// Required component

func NewRelic(
	config *config.Relic,
	logger *zap.Logger,
) *newrelic.Application {
	if config.License == "" {
		logger.Warn("New Relic is not enabled")
		return nil
	}

	app, err := newrelic.NewApplication(
		func(c *newrelic.Config) {
			c.AppName = config.AppName
			c.License = config.License
		},
		nrzap.ConfigLogger(logger),
	)
	if err != nil {
		panic(err)
	}
	return app
}
