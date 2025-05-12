package envx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/envx"
)

// 不存在的环境变量
func TestGetNotExists(t *testing.T) {
	ret := envx.Get("NOT_EXISTS_ENV_KEY", "ENV_VAL")
	assert.Equal(t, "ENV_VAL", ret)
}

// 已存在的环境变量
func TestGetExists(t *testing.T) {
	ret := envx.Get("PATH", "")
	assert.NotEqual(t, "", ret)
}

// 不存在的环境变量
func TestMustGetNotExists(t *testing.T) {
	defer func() {
		assert.Equal(t, "required environment variable NOT_EXISTS_ENV_KEY unset", recover())
	}()

	_ = envx.MustGet("NOT_EXISTS_ENV_KEY")
}

// 已存在的环境变量
func TestMustGetExists(t *testing.T) {
	ret := envx.MustGet("PATH")
	assert.NotEqual(t, "", ret)
}
