package action

import "github.com/gin-gonic/gin"

func (r *router) add(ctx *gin.Context) {
	r.rsp.Ok(ctx, "id", "add")
}
