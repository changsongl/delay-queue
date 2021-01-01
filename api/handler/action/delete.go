package action

import (
	"github.com/gin-gonic/gin"
)

func (r *router) delete(ctx *gin.Context) {
	uriParam := &idTopicParam{}
	err := r.validator.Validate(ctx, uriParam, nil, nil)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	r.rsp.Ok(ctx, uriParam.ID, uriParam.Topic, "delete")
}
