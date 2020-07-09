package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itomofumi/go-gin-xorm-starter/factory"
	"github.com/itomofumi/go-gin-xorm-starter/model"
	"github.com/itomofumi/go-gin-xorm-starter/service"
	"github.com/stretchr/testify/assert"
)

// Setup initializes handler test.
// call this func with "defer".
func Setup() func() {
	gin.SetMode(gin.ReleaseMode)
	// override gin validator
	binding.Validator = &model.StructValidator{}

	return func() {
		gin.SetMode(gin.DebugMode)
	}
}

// ServiceFactoryMock はServiceFactoryのモック実装です
type ServiceFactoryMock struct {
	factory.Servicer
	FruitsMock service.FruitsInterface
	UsersMock  service.UsersInterface
}

// NewFruits returns FruitsMock
func (sf *ServiceFactoryMock) NewFruits() service.FruitsInterface {
	return sf.FruitsMock
}

// NewUsers returns UsersMock
func (sf *ServiceFactoryMock) NewUsers() service.UsersInterface {
	return sf.UsersMock
}

func createGinTestContext(mock *ServiceFactoryMock) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(factory.ServiceKey, mock)
	return c, w
}

func testErrorResponse(t *testing.T, want *model.ErrorResponse, w *httptest.ResponseRecorder) {
	var res *model.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, want, res)
}
