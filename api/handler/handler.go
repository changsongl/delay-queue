package handler

import (
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register() server.RouterFunc
}

func Apply(engine *gin.Engine, handlers ...Handler) {
	for _, h := range handlers {
		regFunc := h.Register()
		regFunc(engine)
	}
}
