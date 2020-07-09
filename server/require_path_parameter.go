package server

import (
	"fmt"
	"strconv"

	"github.com/itomofumi/go-gin-xorm-starter/model"

	"github.com/gin-gonic/gin"
)

// RequirePathParam parses PathParameter as uint64 by given param and then sets it to gin Context.
func RequirePathParam(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param(param), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(400, model.NewErrorResponse("400", model.ErrorParam, fmt.Sprintf("%v must be a positive number", param), err.Error()))
			return
		}
		c.Set(param, id)
		c.Next()
	}
}

// RequireStringPathParam parses PathParameter as string by given param and then sets it to gin Context.
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
