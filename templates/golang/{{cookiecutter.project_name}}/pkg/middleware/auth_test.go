package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/sprig/v3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/account"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/common"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/middleware"
)

func TestUserAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("unittest-session", store))
	router.Use(middleware.UserAuth(account.NewMockAuthBackend()))
	router.SetFuncMap(sprig.FuncMap())
	router.LoadHTMLGlob("../../templates/web/*")

	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get(common.UserIDKey)
		if exists {
			c.String(http.StatusOK, "user id: %s", userID)
		} else {
			c.String(http.StatusInternalServerError, "user id not found")
		}
	})

	t.Run("No user token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid user token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: "InvalidToken"})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Valid user token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: "EverythingIsPermitted"})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "user id: admin", rr.Body.String())
	})
}

func TestUserAuthWithSession(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("unittest-session", store))
	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(common.UserTokenKey, "EverythingIsPermitted")
		session.Set(common.UserIDKey, "session-admin")
		_ = session.Save()
		c.Next()
	})
	router.Use(middleware.UserAuth(account.NewMockAuthBackend()))

	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get(common.UserIDKey)
		if exists {
			c.String(http.StatusOK, "user id: %s", userID)
		} else {
			c.String(http.StatusInternalServerError, "user id not found")
		}
	})

	t.Run("Use username in session", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: "EverythingIsPermitted"})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "user id: session-admin", rr.Body.String())
	})
}
