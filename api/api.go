package api

import (
	"github.com/changsongl/delay-queue/api/handler"
	"github.com/changsongl/delay-queue/api/handler/action"
	"github.com/changsongl/delay-queue/pkg/http"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
	"sync"
)

type API interface {
	RouterFunc() server.RouterFunc
}

type api struct {
	l           log.Logger
	httpHandler handler.Handler
	rsp         http.Response
	sync.Once
}

func NewApi(l log.Logger) API {
	responseHelper := http.Response{}
	httpHandler := action.NewHandler(
		responseHelper, l.WithModule("http"))

	return &api{
		l:           l,
		httpHandler: httpHandler,
		rsp:         responseHelper,
	}
}

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
