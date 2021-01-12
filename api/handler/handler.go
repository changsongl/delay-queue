package handler

import (
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
)

// Handler interface for register all routers for server
type Handler interface {
	Register() server.RouterFunc
}

// Apply register all handler to engine
func Apply(engine *gin.Engine, handlers ...Handler) {
	for _, h := range handlers {
		regFunc := h.Register()
		regFunc(engine)
	}
}
