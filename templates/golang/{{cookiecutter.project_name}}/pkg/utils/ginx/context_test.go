package ginx_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

func TestGetRequestID(t *testing.T) {
	c := &gin.Context{}

	requestID := ginx.GetRequestID(c)
	assert.Equal(t, "", requestID)
}

func TestSetRequestID(t *testing.T) {
	c := &gin.Context{}

	ginx.SetRequestID(c, "test")
	assert.Equal(t, "test", ginx.GetRequestID(c))
}

func TestGetError(t *testing.T) {
	c := &gin.Context{}

	err, ok := ginx.GetError(c)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, err)
}

func TestSetError(t *testing.T) {
	c := &gin.Context{}
	err := errors.New("test")

	ginx.SetError(c, err)
	gErr, ok := ginx.GetError(c)

	assert.Equal(t, true, ok)
	assert.Equal(t, err, gErr)
}

func TestGetUserID(t *testing.T) {
	c := &gin.Context{}

	userID := ginx.GetUserID(c)
	assert.Equal(t, "", userID)
}

func TestSetUserID(t *testing.T) {
	c := &gin.Context{}

	ginx.SetUserID(c, "test")
	assert.Equal(t, "test", ginx.GetUserID(c))
}
