package action

import "github.com/gin-gonic/gin"

func (r *router) finish(ctx *gin.Context) {
	r.rsp.Ok(ctx, "id", "finish")
}
