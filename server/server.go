package server

import (
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/vars"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// RouterFunc is a function resgiter handler to gin.Engine
type RouterFunc func(engine *gin.Engine)

// Server interface for a basic method of http web server,
// which can be ran in few step
type Server interface {
	Init()
	RegisterRouters(regFunc RouterFunc)
	Run(addr string) error
}

// server is Server implementation struct, it has a gin.Engine,
// to save all handlers, logger, events and env.
type server struct {
	r           *gin.Engine
	l           log.Logger
	beforeStart []Event
	afterStop   []Event
	env         vars.Env
}

// Event common event function can be ran before server start,
// or after server stop.
type Event func()

// New return a Server based on Options.
func New(options ...Option) Server {
	s := &server{
		env: vars.EnvRelease,
	}

	for _, opt := range options {
		opt(s)
	}

	if s.env == vars.EnvRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	s.r = r

	return s
}

// Init a server, inject the logger to gin, and register prometheus
// metrics for all apis.
func (s *server) Init() {
	s.r.Use(gin.Recovery())

	if s.l != nil {
		s.r.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.l}))
	} else {
		s.r.Use(gin.LoggerWithConfig(gin.LoggerConfig{}))
	}

	WrapPProf(s.r)

	regMetricFunc := setServerMetricHandlerAndMiddleware()
	regMetricFunc(s.r)
}

// RegisterRouters register router functions for server
func (s *server) RegisterRouters(regFunc RouterFunc) {
	regFunc(s.r)
}

// Run the server, with address. waiting for shutdown signal.
func (s *server) Run(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.r,
	}

	sc := newShutdownChan()

	go func() {
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
		for {
			select {
			case <-term:
				s.l.Info("Signal stop server")
				shutdown(srv, sc)
				return
			}
		}
	}()

	s.l.Info("Run server", log.String("address", addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	if err := sc.Wait(); err != nil {
		return err
	}
	s.l.Info("Server is stopped")

	return nil
}
