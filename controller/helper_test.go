package controller_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup initializes controller test.
// call this func with "defer".
func Setup() func() {
	gin.SetMode(gin.ReleaseMode)
	return func() {
		gin.SetMode(gin.DebugMode)
	}
}

// RegistryMock はServiceRegistryのモック実装です
type RegistryMock struct {
	service.RegistryInterface
	FruitsMock service.FruitsInterface
}

// NewFruits returns FruitsMock
func (r *RegistryMock) NewFruits() service.FruitsInterface {
	return r.FruitsMock
}

func createGinTestContext(registry *RegistryMock) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(service.RegistryKey, registry)
	return c, w
}

func testErrorResponse(t *testing.T, want *model.ErrorResponse, w *httptest.ResponseRecorder) {
	var res *model.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, want, res)
}
