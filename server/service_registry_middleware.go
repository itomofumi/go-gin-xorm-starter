package server

import (
	"github.com/gemcook/go-gin-xorm-starter/service"

	"github.com/gin-gonic/gin"
)

// ServiceRegistryMiddleware provides the service registry
func ServiceRegistryMiddleware(registry service.RegistryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(service.RegistryKey, registry)
		c.Next()
	}
}
