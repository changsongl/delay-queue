package action

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (r *router) add(ctx *gin.Context) {
	uriParam := &topicParam{}
	bodyParam := &addParam{}
	err := r.validator.Validate(ctx, uriParam, nil, bodyParam)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}
	fmt.Printf("%+v", bodyParam)

	r.rsp.Ok(ctx, "id", uriParam.Topic, "value")
}
