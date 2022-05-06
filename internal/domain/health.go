package domain

type HealthService interface {
	IsHealthy() bool
	IsReady() bool

	OnAppStarted()
}
