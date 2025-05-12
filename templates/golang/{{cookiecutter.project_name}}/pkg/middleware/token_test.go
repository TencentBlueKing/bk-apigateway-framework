package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/middleware"
	testingutil "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/utils/testing"
)

func TestQueryTokenAuthRight(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c := testingutil.CreateTestContextWithDefaultRequest(w)

	q := c.Request.URL.Query()
	q.Add("token", "token_for_test")
	c.Request.URL.RawQuery = q.Encode()

	middleware.QueryTokenAuth("token_for_test")(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQueryTokenAuthBad(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c := testingutil.CreateTestContextWithDefaultRequest(w)

	middleware.QueryTokenAuth("token_for_test")(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
