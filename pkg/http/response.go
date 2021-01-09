package http

import (
	"github.com/changsongl/delay-queue/job"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
}

func responseOk(ctx *gin.Context, m map[string]interface{}) {
	ctx.JSON(http.StatusOK, m)
}

func (r *Response) Pong(ctx *gin.Context) {
	responseOk(ctx, map[string]interface{}{"success": true, "value": "pong"})
}

func (r *Response) Ok(ctx *gin.Context) {
	responseOk(ctx, map[string]interface{}{"success": true, "message": "ok"})
}

func (r *Response) OkWithIdAndBody(ctx *gin.Context, id job.Id, value job.Body) {
	responseOk(ctx, map[string]interface{}{"success": true, "id": id, "value": value})
}

func (r *Response) Error(ctx *gin.Context, err error) {
	responseOk(ctx, map[string]interface{}{"success": false, "message": err.Error()})
}
