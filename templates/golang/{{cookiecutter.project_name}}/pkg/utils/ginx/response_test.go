package ginx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

func TestSetSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/success", func(c *gin.Context) {
		ginx.SetResp(c, http.StatusOK, "test data")
	})
	req, _ := http.NewRequest(http.MethodGet, "/success", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	expectedResponse := ginx.Response{Message: "", Data: "test data"}

	var actualResponse ginx.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestSetErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/error", func(c *gin.Context) {
		ginx.SetErrResp(c, http.StatusInternalServerError, "error occurred")
	})
	req, _ := http.NewRequest(http.MethodGet, "/error", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	expectedResponse := ginx.Response{Message: "error occurred", Data: nil}

	var actualResponse ginx.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestNewPaginatedRespData(t *testing.T) {
	data := ginx.NewPaginatedRespData(100, []string{"alpha", "beta", "gamma"})
	assert.Equal(t, ginx.PaginatedResp{Count: int64(100), Results: []string{"alpha", "beta", "gamma"}}, data)
}
