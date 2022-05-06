package domain

import (
	"context"
	"github.com/yuridevx/app-example/api/v1/schedule"
	"time"
)

type Clock interface {
	Now() time.Time
}

type Runner interface {
	Run(ctx context.Context, id string, opts *schedule.JobOptions) error
}
