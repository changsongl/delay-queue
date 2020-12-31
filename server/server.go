package server

import (
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"syscall"
)

type Option func(s *server)

type server struct {
	r           *gin.Engine
	l           log.Logger
	beforeStart []Event
	afterStop   []Event
}

type Event func()

type Server interface {
	Register(r func(s *gin.Engine))
	SetBeforeStartEvent(events ...Event)
	SetAfterStartEvent(events ...Event)
}

func New(options ...Option) Server {
	return server{
		r: gin.Default(),
	}
}

func LoggerOption(l log.Logger) Option {
	return func(s *server) {
		s.l = l
	}
}

func (s server) Register(rf func(r *gin.Engine)) {
	rf(s.r)
}

func (s server) Run(addr string) error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.r,
	}

	sc := NewShutdownChan()

	go func() {
		c := make(chan os.Signal, 1)
		for {
			sig := <-c
			switch sig {
			case syscall.SIGTERM:
				shutdown(srv, sc)
				return
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.l.Fatal(err.Error())
	}

	if err := sc.Wait(); err != nil {
		s.l.Error(err.Error())
	}

	return nil
}

func (s server) SetBeforeStartEvent(events ...Event) {
	s.beforeStart = append(s.beforeStart, events...)
}

func (s server) SetAfterStartEvent(events ...Event) {
	s.beforeStart = append(s.afterStop, events...)
}
