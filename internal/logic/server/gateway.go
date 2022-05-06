package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yuridevx/app"
	"github.com/yuridevx/app-example/api/v1/schedule"
	"github.com/yuridevx/app-example/internal/config"
	"net/http"
)

type Gateway struct {
	config *config.Config
	server *http.Server
	nr     *newrelic.Application
	mux    *runtime.ServeMux
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pattern, ok := runtime.HTTPPathPattern(r.Context())
	var name string
	if ok {
		name = r.Method + " " + pattern
	} else {
		name = r.Method
	}

	txn := g.nr.StartTransaction(name)
	defer txn.End()

	w = txn.SetWebResponse(w)
	txn.SetWebRequestHTTP(r)

	g.mux.ServeHTTP(w, r)
}

func (g *Gateway) ListenGateway(_ context.Context) error {
	err := g.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (g *Gateway) ShutdownGateway(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}

var _ http.Handler = (*Gateway)(nil)

func NewGateway(
	scheduleServer schedule.ScheduleServer,
	config *config.Config,
	nr *newrelic.Application,
	ap app.AppBuilder,
) *Gateway {
	g := &Gateway{
		config: config,
		server: &http.Server{
			Addr: config.ListenHttp,
		},
		nr:  nr,
		mux: runtime.NewServeMux(),
	}
	g.server.Handler = g
	err := schedule.RegisterScheduleHandlerServer(context.Background(), g.mux, scheduleServer)
	if err != nil {
		panic(err)
	}

	ap.C(g).
		PBlocking(g.ListenGateway).
		OnShutdown(g.ShutdownGateway)

	return g
}
