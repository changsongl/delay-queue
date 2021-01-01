package action

import "github.com/gin-gonic/gin"

func (r *router) pop(ctx *gin.Context) {
	r.rsp.Ok(ctx, "id", "pop")
}
