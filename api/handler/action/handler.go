package action

import (
	"github.com/changsongl/delay-queue/api/handler"
	"github.com/changsongl/delay-queue/pkg/http"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
)

type router struct {
	rsp    http.Response
	logger log.Logger
}

func NewHandler(rsp http.Response, logger log.Logger) handler.Handler {
	return &router{
		rsp:    rsp,
		logger: logger,
	}
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
