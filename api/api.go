package api

import (
	"github.com/changsongl/delay-queue/api/handler"
	"github.com/changsongl/delay-queue/api/handler/http"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/pkg/rsp"
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
	sync.Once
}

func NewApi(l log.Logger) API {
	httpHandler := http.NewHandler()

	return &api{
		l:           l,
		httpHandler: httpHandler,
	}
}

func (a *api) RouterFunc() server.RouterFunc {
	rf := func(engine *gin.Engine) {}
	a.Do(func() {
		pingFunc := func(ctx *gin.Context) {
			rsp.ResponsePong(ctx)
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
