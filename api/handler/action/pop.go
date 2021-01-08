package action

import (
	"github.com/gin-gonic/gin"
)

func (r *router) pop(ctx *gin.Context) {
	uriParam := &topicParam{}
	err := r.validator.Validate(ctx, uriParam, nil, nil)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	id, body, err := r.dispatch.Pop(uriParam.Topic)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	r.rsp.OkWithIdAndBody(ctx, id, body)
}
