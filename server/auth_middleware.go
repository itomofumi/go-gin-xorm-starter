package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/util"
	"github.com/gemcook/gognito/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var authContextKey = "auth"

// SetAuth passes an authenticator.
func SetAuth(authenticator auth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(authContextKey, authenticator)
		c.Next()
	}
}

// AuthMiddleware verifies JWT with authenticator.
func AuthMiddleware() gin.HandlerFunc {
	logger := util.GetLogger()
	return func(c *gin.Context) {
		authenticator := c.MustGet(authContextKey).(auth.Authenticator)
		err := authHandler(c, authenticator)
		if err != nil {
			er := model.NewErrorResponse("401", model.ErrorAuth, err.Error())
			logger.Debugln(er)
			c.AbortWithStatusJSON(http.StatusUnauthorized, er)
		}
		c.Next()
	}
}

// OptionalAuthMiddleware does optional JWT verification.
func OptionalAuthMiddleware() gin.HandlerFunc {
	logger := util.GetLogger()
	return func(c *gin.Context) {
		if _, ok := GetBearer(c.Request.Header["Authorization"]); ok {
			authenticator := c.MustGet(authContextKey).(auth.Authenticator)
			err := authHandler(c, authenticator)
			if err != nil {
				er := model.NewErrorResponse("401", model.ErrorAuth, err.Error())
				logger.Debugln(er)
				c.AbortWithStatusJSON(http.StatusUnauthorized, er)
			}
		} else {
			c.Next()
		}
	}
}

func authHandler(c *gin.Context, authenticator auth.Authenticator) error {
	tokenString, ok := GetBearer(c.Request.Header["Authorization"])

	if !ok {
		return fmt.Errorf("Bearer token was not found in Authorization header")
	}

	authedUser, err := authenticateUser(tokenString, authenticator)
	if err != nil {
		return err
	}
	// set user information to Gin's context.
	c.Set("email", authedUser.Email)
	c.Set("sub", authedUser.Sub)
	c.Set("token", authedUser.Token)

	return nil
}

// AuthenticatedUser verified user information
type AuthenticatedUser struct {
	Email string
	Sub   string
	Token *jwt.Token
}

// authenticateUser performs authentication to the given JWT token.
func authenticateUser(tokenString string, authenticator auth.Authenticator) (*AuthenticatedUser, error) {
	token, err := authenticator.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("token is not valid. [Reason] %v", err)
	}

	// "token == nil" may happen when 'NoVerification Mode' is enabled.
	if token == nil || token.Claims == nil {
		return nil, fmt.Errorf("wrong format token")
	}

	// よく使うユーザー情報も渡す
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("wrong format token")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("token must contain email")
	}
	sub := claims["sub"].(string)

	authedUser := AuthenticatedUser{
		Email: email,
		Sub:   sub,
		Token: token,
	}

	return &authedUser, nil
}

// GetBearer gets a bearer token from Authorization header
func GetBearer(auth []string) (jwt string, ok bool) {
	for _, v := range auth {
		ret := strings.Split(v, " ")
		if len(ret) == 2 && ret[0] == "Bearer" {
			return ret[1], true
		}
	}
	return "", false
}
