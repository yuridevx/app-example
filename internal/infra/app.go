package infra

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/yuridevx/app"
	"github.com/yuridevx/app/apprelic"
	"github.com/yuridevx/app/apptrace"
	"github.com/yuridevx/app/appzap"
	"github.com/yuridevx/app/extension"
	"github.com/yuridevx/app/options"
	"go.uber.org/zap"
	"strings"
	"time"
)

func NewApp(
	logger *zap.Logger,
	nr *newrelic.Application,
	pyro *pyroscope.Profiler,
) app.AppBuilder {
	var middleware = []extension.Middleware{
		apptrace.NewTraceMiddleware(
			nil,
		),
		appzap.ZapMiddleware(
			logger.WithOptions(zap.WithCaller(false)),
			func(trace *apptrace.Trace, callType extension.CallType, time appzap.LogTime) (bool, []zap.Field) {
				fields := []zap.Field{
					zap.String("call", strings.Replace(string(callType), "Call", "", -1)),
				}
				if time == appzap.LogBefore && (callType == extension.CallPBlocking ||
					callType == extension.CallStart) {
					return true, fields
				}
				if time == appzap.LogAfter && (callType == extension.CallShutdown) {
					return true, fields
				}
				return trace.GetLog(), fields
			},
		),
		appzap.LogMeMiddleware,
	}
	if nr != nil {
		middleware = append(middleware, apprelic.NewNewRelicMiddleware(nr))
	}
	a := app.NewBuilder(options.ApplicationOptions{
		GlobalMiddleware: middleware,
	})
	a.OnShutdown(func(ctx context.Context) {
		if pyro != nil {
			_ = pyro.Stop()
		}
		_ = logger.Sync()
		nr.Shutdown(time.Second * 15)
	})
	return a
}
