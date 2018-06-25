package controller

import (
	"net/http"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// PostUser は新規ユーザー登録
func PostUser(c *gin.Context) {
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)
	userService := registry.NewUsers()

	body := model.UserCreateBody{}
	err := c.ShouldBindWith(&body, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, "request body mismatch", err.Error()))
		return
	}

	created, err := userService.Create(body.Email, &body.UserProfile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err.Error()))
		return
	}

	c.JSON(http.StatusOK, created)
}
