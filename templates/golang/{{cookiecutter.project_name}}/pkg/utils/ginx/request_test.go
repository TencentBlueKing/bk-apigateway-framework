package ginx_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

func TestGetPage(t *testing.T) {
	tests := []struct {
		name string
		path string
		want int
	}{
		{
			name: "empty",
			path: "/",
			want: 1,
		},
		{
			name: "page=1",
			path: "/?page=1",
			want: 1,
		},
		{
			name: "page=5",
			path: "/?page=5",
			want: 5,
		},
		{
			name: "invalid page",
			path: "/?page=-3",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("GET", tt.path, nil)
			assert.Equal(t, tt.want, ginx.GetPage(c))
		})
	}
}

func TestGetLimit(t *testing.T) {
	tests := []struct {
		name string
		path string
		want int
	}{
		{
			name: "empty",
			path: "/",
			want: 5,
		},
		{
			name: "limit=15",
			path: "/?limit=15",
			want: 15,
		},
		{
			name: "invalid limit",
			path: "/?limit=-1",
			want: 5,
		},
		{
			name: "too large limit",
			path: "/?limit=500",
			want: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("GET", tt.path, nil)
			assert.Equal(t, tt.want, ginx.GetLimit(c))
		})
	}
}

func TestGetOffset(t *testing.T) {
	tests := []struct {
		name string
		path string
		want int
	}{
		{
			name: "empty",
			path: "/",
			want: 0,
		},
		{
			name: "with page and limit",
			path: "/?page=3&limit=15",
			want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("GET", tt.path, nil)
			assert.Equal(t, tt.want, ginx.GetOffset(c))
		})
	}
}
