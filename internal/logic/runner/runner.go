package runner

import (
	"context"
	"github.com/yuridevx/app-example/api/v1/schedule"
	"github.com/yuridevx/app-example/internal/domain"
	"github.com/yuridevx/app/apprelic"
	"go.uber.org/zap"
)

type Runner struct {
	logger *zap.Logger
}

func (r *Runner) Run(ctx context.Context, id string, opts *schedule.JobOptions) error {
	defer apprelic.QuickSegment(ctx)() // Measure time of execution
	r.logger.Info("run job", zap.String("id", id))
	return nil
}

var _ domain.Runner = (*Runner)(nil)

func NewRunner() *Runner {
	return &Runner{}
}
