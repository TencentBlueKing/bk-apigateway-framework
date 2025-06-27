package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/common"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/middleware"
)

func TestAccessControl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allows all users when no configure", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.AccessControl([]string{}))
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "allowed")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "allowed", w.Body.String())
	})

	t.Run("allows specified user", func(t *testing.T) {
		router := gin.New()
		allowedUsers := []string{"user1"}
		router.Use(func(c *gin.Context) {
			c.Set(common.UserIDKey, "user1")
		})
		router.Use(middleware.AccessControl(allowedUsers))
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "allowed")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "allowed", w.Body.String())
	})

	t.Run("forbids unauthorized user", func(t *testing.T) {
		router := gin.New()
		allowedUsers := []string{"user1"}
		router.Use(func(c *gin.Context) {
			c.Set(common.UserIDKey, "user2")
		})
		router.Use(middleware.AccessControl(allowedUsers))
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "forbidden")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
