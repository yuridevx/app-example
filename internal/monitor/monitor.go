package monitor

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yuridevx/app"
	"github.com/yuridevx/app-example/internal/config"
	"github.com/yuridevx/app-example/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

type Monitor struct {
	logger zap.Logger
	server *http.Server
	mux    *http.ServeMux
	health domain.HealthService
}

func (m *Monitor) RunBlocking() {
	m.logger.Info("starting monitor", zap.String("address", m.server.Addr))
	err := m.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		m.logger.Error("unable to start monitoring", zap.Error(err))
		return
	}
	m.logger.Error("monitoring stopped")
}

func (m *Monitor) Stop(ctx context.Context) error {
	err := m.server.Shutdown(ctx)
	return err
}

func (m *Monitor) Health(resp http.ResponseWriter, req *http.Request) {
	if m.health.IsHealthy() {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (m *Monitor) Ready(writer http.ResponseWriter, request *http.Request) {
	if m.health.IsReady() {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func NewMonitor(
	logger zap.Logger,
	ap app.Builder,
	config *config.Config,
	health domain.HealthService,
) *Monitor {
	m := &Monitor{
		logger: logger,
		server: &http.Server{
			Addr: config.ListenMetrics,
		},
		mux:    http.NewServeMux(),
		health: health,
	}
	m.mux.Handle("/metrics", promhttp.Handler())
	m.mux.HandleFunc("/healthz", m.Health)
	m.mux.HandleFunc("/readyz", m.Ready)
	ap.C(m).
		PBlocking(m.RunBlocking). //Run in different goroutine
		OnShutdown(m.Stop)        // Shutdown runs as part of graceful shutdown
	return m
}
