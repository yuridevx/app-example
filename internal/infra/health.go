package infra

import (
	"github.com/yuridevx/app-example/internal/domain"
	"go.uber.org/atomic"
)

type Health struct {
	healthy atomic.Bool
	ready   atomic.Bool
}

func (h *Health) IsHealthy() bool {
	return h.healthy.Load()
}

func (h *Health) IsReady() bool {
	return h.ready.Load()
}

func (h *Health) OnAppStarted() {
	h.ready.Store(true)
	h.healthy.Store(true)
}

var _ domain.HealthService = (*Health)(nil)

func NewHealth() *Health {
	return &Health{}
}
