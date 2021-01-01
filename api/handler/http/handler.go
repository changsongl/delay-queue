package http

import (
	"github.com/changsongl/delay-queue/api/handler"
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
)

type router struct {
}

func NewHandler() handler.Handler {
	return &router{}
}

func (r *router) Register() server.RouterFunc {
	resourcePath := "/job"

	return func(engine *gin.Engine) {
		engine.PUT(resourcePath, r.finish)
		engine.POST(resourcePath, r.add)
		engine.GET(resourcePath, r.pop)
		engine.DELETE(resourcePath, r.delete)
	}
}
