package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Validator interface {
	Validate(c *gin.Context, uri, query, body interface{}) error
}

type paramError struct {
	paramFrom string
	paramErr  error
}

func newParamError(paramFrom string, paramErr error) error {
	return paramError{
		paramFrom: paramFrom,
		paramErr:  paramErr,
	}
}

func (p paramError) Error() string {
	return fmt.Sprintf("%s invalid: %v", p.paramFrom, p.paramErr)
}

type validator struct {
}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(c *gin.Context, uri, query, body interface{}) error {
	if uri != nil {
		if err := c.ShouldBindUri(uri); err != nil {
			return newParamError("uri", err)
		}
	}
	if query != nil {
		if err := c.ShouldBindQuery(query); err != nil {
			return newParamError("query", err)
		}
	}
	if body != nil {
		if err := c.ShouldBindJSON(body); err != nil {
			return newParamError("body", err)
		}
	}

	return nil
}
