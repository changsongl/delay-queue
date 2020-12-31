package server

import (
	"context"
	"net/http"
	"time"
)

const (
	ShutDownTimeout = 6 * time.Second
)

type ShutdownChan chan error

func NewShutdownChan() ShutdownChan {
	return make(chan error, 1)
}

func (s ShutdownChan) Notify(e error) {
	s <- e
}

func (s ShutdownChan) Wait() error {
	return <-s
}

func shutdown(srv *http.Server, sc ShutdownChan) {
	ctx, cancel := context.WithTimeout(context.Background(), ShutDownTimeout)
	sc.Notify(srv.Shutdown(ctx))
	defer cancel()
}
