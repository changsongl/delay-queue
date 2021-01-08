package action

import (
	"github.com/changsongl/delay-queue/type/job"
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

	d, ttr := getDelayAndTTR(bodyParam.Delay, bodyParam.TTR)
	err = r.dispatch.Add(uriParam.Topic, bodyParam.ID, d, ttr, bodyParam.Body)
	if err != nil {
		r.rsp.Error(ctx, err)
		return
	}

	r.rsp.Ok(ctx)
}

func getDelayAndTTR(d, ttr uint) (job.Delay, job.TTR) {
	return job.Delay(d), job.TTR(ttr)
}
