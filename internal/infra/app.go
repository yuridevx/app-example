package infra

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/yuridevx/app"
	"github.com/yuridevx/app/apprelic"
	"github.com/yuridevx/app/apptrace"
	"github.com/yuridevx/app/appzap"
	"github.com/yuridevx/app/options"
	"go.uber.org/zap"
	"strings"
	"time"
)

func NewApp(
	logger *zap.Logger,
	nr *newrelic.Application,
	pyro *pyroscope.Profiler,
) app.Builder {
	var middleware = []options.Middleware{
		apptrace.NewTraceMiddleware(
			nil,
			nil,
		),
		appzap.ZapMiddleware(
			logger.WithOptions(zap.WithCaller(false)),
			func(trace *apptrace.Trace, call options.Call, time appzap.LogTime) (bool, []zap.Field) {
				callType := call.GetCallType()
				fields := []zap.Field{
					zap.String("call", strings.Replace(string(callType), "Call", "", -1)),
				}
				if time == appzap.LogBefore && (callType == options.CallPBlocking ||
					callType == options.CallStart) {
					return true, fields
				}
				if time == appzap.LogAfter && (callType == options.CallShutdown) {
					return true, fields
				}
				return trace.GetLog(), fields
			},
		),
		appzap.LogMeMiddleware,
	}
	if nr != nil {
		middleware = append(middleware,
			apprelic.NewRelicTransactionMiddleware(nr),
			apprelic.NewRelicTraceMiddleware(),
		)
	}
	a := app.NewBuilder(func(o *options.ApplicationOptions) {
		o.Middleware = middleware
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
