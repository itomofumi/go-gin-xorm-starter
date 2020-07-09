package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/itomofumi/go-gin-xorm-starter/server"
	"github.com/itomofumi/go-gin-xorm-starter/util"
	"github.com/stretchr/testify/assert"
)

func TestLogMiddleware(t *testing.T) {
	called := false

	util.GetTimeNowFunc = func() time.Time {
		t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		return t
	}

	router := gin.Default()
	b := &bytes.Buffer{}
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Out = b
	logger.Formatter = &logrus.JSONFormatter{}

	router.Use(server.LogMiddleware(logger, time.RFC3339, false))
	router.GET("/v1/tests", func(c *gin.Context) {
		called = true
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/tests?param=123", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("User-Agent", "httptest")
	router.ServeHTTP(w, req)

	assert := assert.New(t)
	assert.True(called)

	j := map[string]interface{}{}
	_ = json.Unmarshal(b.Bytes(), &j)
	t.Log(j)
	assert.Equal("info", j["level"])
	assert.Equal("[access]", j["msg"])
	assert.Equal("GET", j["method"])
	assert.Equal(float64(200), j["status"])
	assert.Equal("/v1/tests", j["path"])
	assert.Equal("param=123", j["query"])
	assert.Equal([]interface{}{"https://example.com"}, j["origin"])
	assert.Equal("httptest", j["user-agent"])
	assert.Equal("2009-11-10T23:00:00Z", j["fields.time"])
}
