package server

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware appends CORS headers
func CORSMiddleware() gin.HandlerFunc {
	var accessControlAllowHeaders = strings.Join(
		[]string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
		", ",
	)

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, UPDATE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", accessControlAllowHeaders)
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "GET" {
			c.Writer.Header().Set("Access-Control-Max-Age", "0")
		} else {
			c.Writer.Header().Set("Access-Control-Max-Age", "0")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
