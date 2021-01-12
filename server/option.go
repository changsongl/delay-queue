package server

import (
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/vars"
)

// Option is a function to the server
// for setting options
type Option func(s *server)

// LoggerOption set logger to server, server uses it
// to log result of requests
func LoggerOption(l log.Logger) Option {
	return func(s *server) {
		s.l = l.WithModule("server")
	}
}

// EnvOption set env for server.
func EnvOption(env vars.Env) Option {
	return func(s *server) {
		s.env = env
	}
}

// BeforeStartEventOption set before start events.
func BeforeStartEventOption(events ...Event) Option {
	return func(s *server) {
		s.beforeStart = append(s.beforeStart, events...)
	}
}

// AfterStopEventOption set after stop events.
func AfterStopEventOption(events ...Event) Option {
	return func(s *server) {
		s.beforeStart = append(s.afterStop, events...)
	}
}
