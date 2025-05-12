package uuidx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/uuidx"
)

func TestNew(t *testing.T) {
	assert.Equal(t, 32, len(uuidx.New()))
}
