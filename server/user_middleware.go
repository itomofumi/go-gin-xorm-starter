package server

import (
	"fmt"
	"net/http"

	"github.com/gemcook/go-gin-xorm-starter/factory"
	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/util"
	"github.com/gemcook/ptr"

	"github.com/gin-gonic/gin"
)

// UserMiddleware 認証したユーザー情報を取得する
func UserMiddleware() gin.HandlerFunc {
	logger := util.GetLogger()
	return func(c *gin.Context) {
		err := UserHandler(c)
		if err != nil {
			er := model.NewErrorResponse("401", model.ErrorAuth, err.Error())
			logger.Debug(er)
			c.AbortWithStatusJSON(http.StatusUnauthorized, er)
			return
		}

		c.Next()
	}
}

// OptionalUserMiddleware 認証していればユーザー情報を取得する
func OptionalUserMiddleware() gin.HandlerFunc {
	logger := util.GetLogger()
	return func(c *gin.Context) {
		_, ok := c.Get("email")
		if ok {
			err := UserHandler(c)
			if err != nil {
				er := model.NewErrorResponse("401", model.ErrorAuth, err.Error())
				logger.Debug(er)
				c.AbortWithStatusJSON(http.StatusUnauthorized, er)
				return
			}
		}

		c.Next()
	}
}

// UserHandler は認証情報からユーザー取得を行う
func UserHandler(c *gin.Context) error {
	factory := c.MustGet(factory.ServiceKey).(factory.Servicer)
	userSrv := factory.NewUsers()

	// ID Tokenのemailはエンドユーザー身元識別子
	email := c.MustGet("email").(string)
	user, ok := userSrv.GetByEmail(email)
	if !ok {
		return fmt.Errorf("cannot find user email = %v", email)
	}

	// ここまで来たらユーザー認証OK
	if user.EmailVerified != nil && !*user.EmailVerified {
		err := userSrv.Verify(user.ID)
		if err != nil {
			return err
		}
		user.EmailVerified = ptr.Bool(true)
	}

	// PublicDataの更新
	user.UserPublicData = *user.GetPublicData()
	c.Set("user", user)
	return nil
}
