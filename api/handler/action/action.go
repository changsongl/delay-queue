package action

import (
	"github.com/changsongl/delay-queue/api/handler"
	"github.com/changsongl/delay-queue/dispatch"
	"github.com/changsongl/delay-queue/job"
	"github.com/changsongl/delay-queue/pkg/http"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/server"
	"github.com/gin-gonic/gin"
	"time"
)

type idParam struct {
	ID job.ID `uri:"id" json:"id" binding:"required,max=200"`
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
	Delay    uint     `json:"delay"`
	TTR      uint     `json:"ttr"`
	Body     job.Body `json:"body"`
	Override bool     `json:"override"`
}

// router container all api actions
type router struct {
	rsp       http.Response
	logger    log.Logger
	validator http.Validator
	dispatch  dispatch.Dispatch
}

// NewHandler return a router
func NewHandler(rsp http.Response, logger log.Logger, validator http.Validator, dispatch dispatch.Dispatch) handler.Handler {
	return &router{
		rsp:       rsp,
		logger:    logger.WithModule("handler"),
		validator: validator,
		dispatch:  dispatch,
	}
}

// Register return a register function for all routers
func (r *router) Register() server.RouterFunc {
	return func(engine *gin.Engine) {
		engine.PUT("/topic/:topic/job/:id", r.finish)
		engine.POST("/topic/:topic/job", r.add)
		engine.GET("/topic/:topic/job", r.pop)
		engine.DELETE("/topic/:topic/job/:id", r.delete)
	}
}

// add action is to push job to delay queue
func (r *router) add(ctx *gin.Context) {
	uriParam := &topicParam{}
	bodyParam := &addParam{}
	err := r.validator.Validate(ctx, uriParam, nil, bodyParam)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	d, ttr := getDelayAndTTR(bodyParam.Delay, bodyParam.TTR)
	err = r.dispatch.Add(uriParam.Topic, bodyParam.ID, d, ttr, bodyParam.Body, bodyParam.Override)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	r.rsp.Ok(ctx)
}

// getDelayAndTTR convert user seconds to delay object
func getDelayAndTTR(d, ttr uint) (job.Delay, job.TTR) {
	second := uint(time.Second)
	return job.Delay(d * second), job.TTR(ttr * second)
}

// delete is for deleting job for running
func (r *router) delete(ctx *gin.Context) {
	uriParam := &idTopicParam{}
	err := r.validator.Validate(ctx, uriParam, nil, nil)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	err = r.dispatch.Delete(uriParam.Topic, uriParam.ID)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	r.rsp.Ok(ctx)
}

// finish is for ack job, which is just processed by the user.
// it means delay queue won't retry to send this job to user
// again.
func (r *router) finish(ctx *gin.Context) {
	uriParam := &idTopicParam{}
	err := r.validator.Validate(ctx, uriParam, nil, nil)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	err = r.dispatch.Finish(uriParam.Topic, uriParam.ID)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	r.rsp.Ok(ctx)
}

// pop a job from delay queue, if there is no job in the ready queue,
// then id and topic are empty
func (r *router) pop(ctx *gin.Context) {
	uriParam := &topicParam{}
	err := r.validator.Validate(ctx, uriParam, nil, nil)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	j, err := r.dispatch.Pop(uriParam.Topic)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	r.rsp.OkWithJob(ctx, j)
}
