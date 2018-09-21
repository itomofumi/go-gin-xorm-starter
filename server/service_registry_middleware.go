package server

import (
	"github.com/gemcook/go-gin-xorm-starter/factory"
	"github.com/gin-gonic/gin"
)

// ServiceKeyMiddleware provides the service factory
func ServiceKeyMiddleware(si factory.ServiceInitializer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(factory.ServiceKey, si)
		c.Next()
	}
}
