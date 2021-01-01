package http

import (
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

func (r *Response) Error(ctx *gin.Context, err error) {
	responseOk(ctx, map[string]interface{}{"success": false, "error": err.Error()})
}

func (r *Response) Ok(ctx *gin.Context, id, topic, value string) {
	responseOk(ctx, map[string]interface{}{"success": true, "id": id, "topic": topic, "value": value})
}
