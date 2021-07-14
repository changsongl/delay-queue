package http

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	SuccessMessage = "ok"
)

type Response struct {
}

func responseOk(ctx *gin.Context, m map[string]interface{}) {
	ctx.JSON(http.StatusOK, m)
}

func (r *Response) Pong(ctx *gin.Context) {
	responseOk(ctx, map[string]interface{}{"success": true, "message": "pong"})
}

func (r *Response) Ok(ctx *gin.Context) {
	responseOk(ctx, map[string]interface{}{
		"success": true,
		"message": SuccessMessage,
	})
}

func (r *Response) OkWithJob(ctx *gin.Context, j *job.Job) {
	var jobMap map[string]interface{}
	if j != nil {
		jobMap = map[string]interface{}{
			"topic": j.Topic,
			"id":    j.ID,
			"body":  j.Body,
			"ttr":   time.Duration(j.TTR) / time.Second,
			"delay": time.Duration(j.Delay) / time.Second,
		}
	}

	responseOk(ctx, map[string]interface{}{
		"success": true,
		"message": SuccessMessage,
		"data":    jobMap,
	})
}

func (r *Response) Error(ctx *gin.Context, err error) {
	responseOk(ctx, map[string]interface{}{
		"success": false,
		"message": err.Error(),
	})
}
