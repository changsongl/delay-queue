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

type Option func(s *server)

type server struct {
	r           *gin.Engine
	l           log.Logger
	beforeStart []Event
	afterStop   []Event
	env         vars.Env
}

type Event func()

type Server interface {
	Init(rfs ...func(r *gin.Engine))
	SetBeforeStartEvent(events ...Event)
	SetAfterStartEvent(events ...Event)
	Run(addr string) error
}

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

func LoggerOption(l log.Logger) Option {
	return func(s *server) {
		s.l = l
	}
}

func EnvOption(env vars.Env) Option {
	return func(s *server) {
		s.env = env
	}
}

func (s *server) Init(rfs ...func(r *gin.Engine)) {
	s.r.Use(gin.Recovery())

	if s.l != nil {
		s.r.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.l}))
	} else {
		s.r.Use(gin.LoggerWithConfig(gin.LoggerConfig{}))
	}

	regMetricFunc := getServerMetricRegisterFunc()
	regMetricFunc(s.r)

	for _, rf := range rfs {
		rf(s.r)
	}
}

func (s *server) Run(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.r,
	}

	sc := NewShutdownChan()

	go func() {
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		for {
			select {
			case <-term:
				shutdown(srv, sc)
				return
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	if err := sc.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *server) SetBeforeStartEvent(events ...Event) {
	s.beforeStart = append(s.beforeStart, events...)
}

func (s *server) SetAfterStartEvent(events ...Event) {
	s.beforeStart = append(s.afterStop, events...)
}
