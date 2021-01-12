package api

import (
	"github.com/changsongl/delay-queue/api/handler"
	"github.com/changsongl/delay-queue/api/handler/action"
	"github.com/changsongl/delay-queue/dispatch"
	"github.com/changsongl/delay-queue/pkg/http"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
	"sync"
)

// API interface, return a router func to register
// all handlers for gin engine
type API interface {
	RouterFunc() server.RouterFunc
}

// api struct implemented API
type api struct {
	l           log.Logger
	httpHandler handler.Handler
	rsp         http.Response
	sync.Once
}

// NewApi with logger object and dispatch
func NewApi(l log.Logger, dispatch dispatch.Dispatch) API {
	logger := l.WithModule("api")
	responseHelper := http.Response{}
	httpHandler := action.NewHandler(
		responseHelper, logger, http.NewValidator(), dispatch)

	return &api{
		l:           logger,
		httpHandler: httpHandler,
		rsp:         responseHelper,
	}
}

// RouterFunc return a server.RouterFunc which register
// ping and all delay queue handler actions.
func (a *api) RouterFunc() server.RouterFunc {
	rf := func(engine *gin.Engine) {}

	a.Do(func() {
		pingFunc := func(ctx *gin.Context) {
			a.rsp.Pong(ctx)
		}

		rf = func(engine *gin.Engine) {
			engine.GET("ping", pingFunc)

			handler.Apply(
				engine,
				a.httpHandler,
			)
		}
	})

	return rf
}
