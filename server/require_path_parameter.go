package server

import (
	"fmt"
	"strconv"

	"github.com/gemcook/go-gin-xorm-starter/model"

	"github.com/gin-gonic/gin"
)

// RequirePathParam はPathParameterにparamを要求します
func RequirePathParam(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param(param), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(400, model.NewErrorResponse("400", model.ErrorParam, fmt.Sprintf("%v must be a positive number", param), err.Error()))
			return
		}
		c.Set(param, id)
		c.Next()
	}
}

// RequireStringPathParam はPathParameterにString型のparamを要求します
func RequireStringPathParam(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(param)
		if id == "" {
			c.AbortWithStatusJSON(400, model.NewErrorResponse("400", model.ErrorParam, fmt.Sprintf("%v must not be empty", param)))
			return
		}
		c.Set(param, id)
		c.Next()
	}
}
