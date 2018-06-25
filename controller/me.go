package controller

import (
	"fmt"
	"net/http"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/service"

	"github.com/gin-gonic/gin"
)

// GetMe はログイン情報を取得します
func GetMe(c *gin.Context) {
	registry := c.MustGet(service.RegistryKey).(service.RegistryInterface)
	usersService := registry.NewUsers()

	email := c.MustGet("email").(string)

	user, ok := usersService.GetByEmail(email)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, fmt.Errorf("user not found")))
		return
	}

	c.JSON(http.StatusOK, user)
}
