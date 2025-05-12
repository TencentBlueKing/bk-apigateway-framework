package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/middleware"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

func TestRequestIDMiddlewareWithoutRequestID(t *testing.T) {
	t.Parallel()

	// request with no request_id
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/ping", func(c *gin.Context) {
		requestID := ginx.GetRequestID(c)
		assert.NotNil(t, requestID)
		assert.Len(t, requestID, 32)
		c.String(http.StatusOK, "pong")
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}

func TestRequestIDMiddlewareWithRequestID(t *testing.T) {
	t.Parallel()

	// request with X-Request-Id
	req, _ := http.NewRequest("GET", "/ping", nil)
	originRID := "ca7ff4ce433447a99e8175f28af31460"
	req.Header.Set(common.RequestIDHeaderKey, originRID)

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/ping2", func(c *gin.Context) {
		requestID := ginx.GetRequestID(c)
		assert.NotNil(t, requestID)
		assert.Equal(t, originRID, requestID)
		c.String(http.StatusOK, "pong")
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}
