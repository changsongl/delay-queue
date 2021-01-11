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

type RouterFunc func(engine *gin.Engine)

type server struct {
	r           *gin.Engine
	l           log.Logger
	beforeStart []Event
	afterStop   []Event
	env         vars.Env
}

type Event func()

type Server interface {
	Init()
	RegisterRouters(regFunc RouterFunc)
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

func (s *server) Init() {
	s.r.Use(gin.Recovery())

	if s.l != nil {
		s.r.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: s.l}))
	} else {
		s.r.Use(gin.LoggerWithConfig(gin.LoggerConfig{}))
	}

	regMetricFunc := getServerMetricRegisterFunc()
	regMetricFunc(s.r)
}

func (s *server) RegisterRouters(regFunc RouterFunc) {
	regFunc(s.r)
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
