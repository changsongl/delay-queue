package action

import "github.com/gin-gonic/gin"

func (r *router) delete(ctx *gin.Context) {
	r.rsp.Ok(ctx, "id", "delete")
}
