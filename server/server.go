package server

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/crypto-pe/api"
	"github.com/crypto-pe/api/config"
	"github.com/crypto-pe/api/data"
	"github.com/crypto-pe/api/rpc"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	Config *config.Config
	Logger zerolog.Logger
	RPC    *rpc.RPC

	ctx       context.Context
	ctxStopFn context.CancelFunc
	running   int32
}

func NewServer(cfg *config.Config) (*Server, error) {
	var err error

	cfg.GitCommit = api.GITCOMMIT

	// Logging
	logger := httplog.NewLogger("api", httplog.Options{
		LogLevel:       cfg.Logging.Level,
		LevelFieldName: "severity",
		JSON:           cfg.Logging.JSON,
		Concise:        cfg.Logging.Concise,
		Tags: map[string]string{
			"serviceName":    cfg.Service.Name,
			"serviceVersion": api.GITCOMMIT,
		},
	})

	//
	// Database
	//
	err = data.PrepareDB(cfg)
	if err != nil {
		return nil, err
	}

	// WebRPC Server
	rpc, err := rpc.NewRPC(cfg, logger)
	if err != nil {
		return nil, err
	}

	//
	// Server
	//
	server := &Server{
		Config: cfg,
		Logger: logger,
		RPC:    rpc,
	}

	return server, nil
}

func (s *Server) Run() error {
	if s.IsRunning() {
		return fmt.Errorf("server already running")
	}

	defer s.Stop()

	oplog := s.Logger.With().Str("op", "run").Logger()
	oplog.Info().Msgf("=> run service")

	// Running
	atomic.StoreInt32(&s.running, 1)

	// Server root context
	s.ctx, s.ctxStopFn = context.WithCancel(context.Background())

	// Subprocess run context
	g, ctx := errgroup.WithContext(s.ctx)

	// RPC
	g.Go(func() error {
		oplog.Info().Msgf("-> rpc: run")
		return s.RPC.Run(ctx)
	})

	// Once run context is done, trigger a server-stop.
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	// Wait for subprocesses to finish
	return g.Wait()
}

func (s *Server) Stop() {
	if !s.IsRunning() || s.IsStopping() {
		return
	}

	oplog := s.Logger.With().Str("op", "shutdown").Logger()
	oplog.Info().Msg("=> shutdown service")

	// Stopping
	atomic.StoreInt32(&s.running, 2)

	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.RPC.Stop(shutdownCtx)
	}()

	// Force shutdown after grace period
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			s.Fatal("graceful shutdown timed out.. forced exit.")
		}
	}()

	// Wait for subprocesses to gracefully stop
	wg.Wait()
	s.ctxStopFn()
	atomic.StoreInt32(&s.running, 0)
}

func (s *Server) IsRunning() bool {
	return atomic.LoadInt32(&s.running) >= 1
}

func (s *Server) IsStopping() bool {
	return atomic.LoadInt32(&s.running) == 2
}

func (s *Server) Fatal(format string, v ...interface{}) {
	s.Logger.Fatal().Msgf(format, v...)
}

func (s *Server) End() {
	s.Logger.Info().Msgf("bye.")
}
