package server

import (
	"context"
	"net/http"
	"time"
)

const (
	// shutdown timeout for server, it time comes, and
	// server is not quit yet. it will force server to stop.
	shutDownTimeout = 6 * time.Second
)

// shutdownChan shutdown channel with error
type shutdownChan chan error

// return a shutdownChan
func newShutdownChan() shutdownChan {
	return make(chan error, 1)
}

// Notify the shutdown channel to stop the server
func (s shutdownChan) Notify(e error) {
	s <- e
}

// Wait for shutdown finished
func (s shutdownChan) Wait() error {
	return <-s
}

// shutdown use shutDownTimeout context to notify the shutdownChan
func shutdown(srv *http.Server, sc shutdownChan) {
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
	defer cancel()
	sc.Notify(srv.Shutdown(ctx))
}
