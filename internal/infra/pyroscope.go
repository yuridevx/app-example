package infra

import (
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/yuridevx/app-example/internal/config"
	"go.uber.org/zap"
	"os"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func (z *ZapLogger) Infof(format string, args ...interface{}) {
	z.sugar.Infof(format, args...)
}

func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	z.sugar.Debugf(format, args...)
}

func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	z.sugar.Errorf(format, args...)
}

// Optional component

func NewPyroscope(
	logger *zap.Logger,
	config *config.Pyroscope,
) *pyroscope.Profiler {
	if config.ServerAddress == "" {
		logger.Warn("Pyroscope server address not set, disabling")
		return nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	app, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: config.AppPrefix + hostname,
		ServerAddress:   config.ServerAddress,
		Logger: &ZapLogger{
			sugar: logger.Sugar(),
		},
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})
	if err != nil {
		panic(err)
	}
	return app
}
