package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itomofumi/go-gin-xorm-starter/factory"
	"github.com/itomofumi/go-gin-xorm-starter/model"
)

// GetMe はログイン情報を取得します
func GetMe(c *gin.Context) {
	factory := c.MustGet(factory.ServiceKey).(factory.Servicer)
	usersService := factory.NewUsers()

	email := c.MustGet("email").(string)

	user, ok := usersService.GetByEmail(email)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, fmt.Errorf("user not found")))
		return
	}

	c.JSON(http.StatusOK, user)
}

// PostUser は新規ユーザー登録
func PostUser(c *gin.Context) {
	factory := c.MustGet(factory.ServiceKey).(factory.Servicer)
	userService := factory.NewUsers()

	body := model.UserCreateBody{}
	err := c.ShouldBindWith(&body, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, "request body mismatch", err))
		return
	}

	created, err := userService.Create(body.Email, &body.UserProfile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse("400", model.ErrorParam, err))
		return
	}

	c.JSON(http.StatusCreated, created)
}
