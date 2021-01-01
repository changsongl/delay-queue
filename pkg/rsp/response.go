package rsp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func responseOk(ctx *gin.Context, m map[string]interface{}) {
	ctx.JSON(http.StatusOK, m)
}

func ResponsePong(ctx *gin.Context) {
	responseOk(ctx, map[string]interface{}{"command": "ping", "body": "pong"})
}
