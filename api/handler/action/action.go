package action

import (
	"github.com/changsongl/delay-queue/api/handler"
	"github.com/changsongl/delay-queue/dispatch"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/http"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
)

type idParam struct {
	ID job.Id `uri:"id" json:"id" binding:"required,max=200"`
}

type topicParam struct {
	Topic job.Topic `uri:"topic" json:"topic" binding:"required,max=50"`
}

type idTopicParam struct {
	idParam
	topicParam
}

type addParam struct {
	idParam
	Delay uint     `json:"delay"`
	TTR   uint     `json:"ttr"`
	Body  job.Body `json:"body"`
}

type router struct {
	rsp       http.Response
	logger    log.Logger
	validator http.Validator
	dispatch  dispatch.Dispatch
}

func NewHandler(rsp http.Response, logger log.Logger, validator http.Validator, dispatch dispatch.Dispatch) handler.Handler {
	return &router{
		rsp:       rsp,
		logger:    logger.WithModule("handler"),
		validator: validator,
		dispatch:  dispatch,
	}
}

func (r *router) Register() server.RouterFunc {
	return func(engine *gin.Engine) {
		engine.PUT("/topic/:topic/job/:id", r.finish)
		engine.POST("/topic/:topic/job", r.add)
		engine.GET("/topic/:topic/job", r.pop)
		engine.DELETE("/topic/:topic/job/:id", r.delete)
	}
}
