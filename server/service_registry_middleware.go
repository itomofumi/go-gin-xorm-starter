package server

import (
	"github.com/gin-gonic/gin"
	"github.com/itomofumi/go-gin-xorm-starter/factory"
)

// ServiceKeyMiddleware provides the service factory
func ServiceKeyMiddleware(si factory.Servicer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(factory.ServiceKey, si)
		c.Next()
	}
}
