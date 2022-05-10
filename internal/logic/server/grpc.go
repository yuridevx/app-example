package server

import (
	"context"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yuridevx/app"
	"github.com/yuridevx/app-example/api/v1/schedule"
	"github.com/yuridevx/app-example/internal/config"
	"google.golang.org/grpc"
	"net"
)

type Grpc struct {
	conf   *config.Config
	grpc   *grpc.Server
	listen net.Listener
}

func (s *Grpc) ListenGrpc(_ context.Context) error {
	lis, err := net.Listen("tcp", s.conf.ListenGrpc)
	if err != nil {
		return err
	}
	s.listen = lis
	return nil
}

func (s *Grpc) ServeGrpc(err error) error {
	err = s.grpc.Serve(s.listen)
	if err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

func (s *Grpc) ShutdownGrpc() {
	// It's fine if it takes time.
	// Every component shutdown is run in parallel
	s.grpc.GracefulStop()
}

func NewGrpc(
	conf *config.Config,
	ap app.Builder,
	nr *newrelic.Application,
	scheduleServer schedule.ScheduleServer,
) *Grpc {
	s := &Grpc{
		conf: conf,
		grpc: grpc.NewServer(
			grpc.UnaryInterceptor(nrgrpc.UnaryServerInterceptor(nr)),
			grpc.StreamInterceptor(nrgrpc.StreamServerInterceptor(nr)),
		),
	}

	schedule.RegisterScheduleServer(s.grpc, scheduleServer)
	// If start fails the whole app will fail
	// Every other function is guaranteed to run after OnStart
	// completes successfully
	ap.C(s).
		OnStart(s.ListenGrpc).
		PBlocking(s.ServeGrpc).
		OnShutdown(s.ShutdownGrpc)

	return s
}
