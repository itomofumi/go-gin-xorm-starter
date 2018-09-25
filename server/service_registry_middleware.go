package server

import (
	"github.com/gemcook/go-gin-xorm-starter/factory"
	"github.com/gin-gonic/gin"
)

// ServiceKeyMiddleware provides the service factory
func ServiceKeyMiddleware(si factory.Servicer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(factory.ServiceKey, si)
		c.Next()
	}
}
