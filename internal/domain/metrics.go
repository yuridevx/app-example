package domain

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const promPrefix = "example_"

var JobScheduledCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: promPrefix + "job_scheduled_total",
		Help: "Total number of jobs scheduled",
	},
)

var JobFiredCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: promPrefix + "job_fired_total",
		Help: "Total number of jobs fired",
	},
)

var JobCompleteCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: promPrefix + "job_complete_total",
		Help: "Total number of jobs completed",
	},
)

var JobRandomlyRemoved = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: promPrefix + "job_randomly_removed_total",
		Help: "Total number of jobs removed randomly",
	},
)
