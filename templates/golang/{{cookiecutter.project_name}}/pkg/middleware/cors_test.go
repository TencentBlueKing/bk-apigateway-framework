package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/middleware"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.CORS([]string{"http://example.com", "http://example.org"}))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name               string
		origin             string
		expectedStatusCode int
	}{
		{
			"Allowed origin",
			"http://example.com",
			http.StatusOK,
		},
		{
			"Not allowed origin",
			"http://notallowed.com",
			http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Origin", tt.origin)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
