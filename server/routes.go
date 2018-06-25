package server

import (
	"github.com/gemcook/go-gin-xorm-starter/controller"

	"github.com/gin-gonic/gin"
)

func defineRoutes(r gin.IRouter) {

	// v1
	v1 := r.Group("/v1")
	v1withUser := v1.Group("/", AuthMiddleware(), UserMiddleware())

	{
		v1withUser.GET("/me", controller.GetMe)
		v1withUser.GET("/user", controller.GetMe)
	}

	{
		v1.POST("/users", controller.PostUser)
	}

	{
		v1.GET("/fruits", controller.GetFruits)
		v1.GET("/fruits/:fruit-id", RequirePathParam("fruit-id"), controller.GetFruitByID)
		v1withUser.POST("/fruits", controller.PostFruit)
		v1withUser.PUT("/fruits/:fruit-id", RequirePathParam("fruit-id"), controller.PutFruit)
		v1withUser.DELETE("/fruits/:fruit-id", RequirePathParam("fruit-id"), controller.DeleteFruit)
	}
}
