package server

import (
	"net/http"

	"github.com/itomofumi/go-gin-xorm-starter/handler"

	"github.com/gin-gonic/gin"
)

func defineRoutes(r gin.IRouter) {

	// Health Check
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// v1
	v1 := r.Group("/v1")
	v1withUser := v1.Group("/", AuthMiddleware(), UserMiddleware())

	{
		v1withUser.GET("/me", handler.GetMe)
		v1withUser.GET("/user", handler.GetMe)
	}

	{
		v1.POST("/users", handler.PostUser)
	}

	{
		v1.GET("/fruits", handler.GetFruits)
		v1.GET("/fruits/:fruit-id", RequirePathParam("fruit-id"), handler.GetFruitByID)
		v1withUser.POST("/fruits", handler.PostFruit)
		v1withUser.PUT("/fruits/:fruit-id", RequirePathParam("fruit-id"), handler.PutFruit)
		v1withUser.DELETE("/fruits/:fruit-id", RequirePathParam("fruit-id"), handler.DeleteFruit)
	}
}
