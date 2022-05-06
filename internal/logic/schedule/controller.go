package schedule

import (
	"context"
	"github.com/yuridevx/app"
	"github.com/yuridevx/app-example/api/v1/schedule"
	"github.com/yuridevx/app-example/internal/domain"
	"github.com/yuridevx/app/apperr"
	"github.com/yuridevx/app/apptrace"
	"go.uber.org/multierr"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var ScheduleError = apperr.Error("schedule error")

type Controller struct {
	clock  domain.Clock
	jobs   map[string]*schedule.JobOptions
	runAt  map[string]time.Time
	putCh  chan *schedule.JobRequest
	runner domain.Runner
}

func (c *Controller) NewJob(ctx context.Context, request *schedule.JobRequest) (*schedule.JobResponse, error) {
	if request.Id == "" {
		request.Id = strconv.Itoa(rand.Intn(math.MaxInt))
	}
	select {
	case c.putCh <- request:
		domain.JobScheduledCounter.Inc()
		return &schedule.JobResponse{}, nil
	case <-ctx.Done():
		return nil, ScheduleError(ctx.Err())
	}
}

func (c *Controller) RandomCleanup(ctx context.Context) error {
	if rand.Intn(100) <= 30 {
		return ScheduleError("unlucky random action")
	}
	cnt := 0
	for key, _ := range c.jobs {
		if rand.Intn(100) <= 30 {
			cnt++
			delete(c.jobs, key)
			delete(c.runAt, key)
		}
	}
	domain.JobRandomlyRemoved.Add(float64(cnt))
	if cnt > 0 {
		apptrace.Attributes(ctx, "removed", cnt)
	} else {
		apptrace.DontLog(ctx)
	}
	return nil
}

func (c *Controller) ScheduleTick(ctx context.Context) error {
	cnt := 0
	now := c.clock.Now()
	var err error
	for k, v := range c.runAt {
		if now.After(v) {
			err = multierr.Append(
				err,
				c.runner.Run(ctx, k, c.jobs[k]),
			)
			cnt++
			delete(c.runAt, k)
			delete(c.jobs, k)
		}
	}
	domain.JobFiredCounter.Add(float64(cnt))
	apptrace.Attributes(ctx, "runCounter", cnt)
	return err
}

func (c *Controller) Put(ctx context.Context, req *schedule.JobRequest) error {
	dur, err := time.ParseDuration(req.Delay)
	if err != nil {
		return err
	}
	c.jobs[req.Id] = req.Options
	c.runAt[req.Id] = c.clock.Now().Add(dur)
	apptrace.Attributes(ctx, "scheduled", req.Id)
	return nil
}

var _ schedule.ScheduleServer = (*Controller)(nil)

func NewController(
	ap app.AppBuilder,
	clock domain.Clock,
	runner domain.Runner,
) *Controller {
	c := &Controller{
		runner: runner,
		clock:  clock,
		jobs:   make(map[string]*schedule.JobOptions),
		runAt:  make(map[string]time.Time),
		putCh:  make(chan *schedule.JobRequest, 100),
	}
	// We don't need to use any kind of synchronization here,
	// because Compete Consume / Compete Period / Compete guarantees that
	// all functions will be executed on the same goroutine.
	ap.C(c).
		CConsume(c.putCh, c.Put).
		CPeriod(time.Second*10, c.ScheduleTick).
		CPeriod(time.Second*5, c.RandomCleanup)

	return c
}
