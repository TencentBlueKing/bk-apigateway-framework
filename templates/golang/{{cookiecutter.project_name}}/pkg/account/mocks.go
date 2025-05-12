package account

import (
	"context"

	"github.com/pkg/errors"
)

// MockAuthBackend 测试用 AuthBackend
type MockAuthBackend struct{}

// GetLoginUrl 获取登录地址
func (b *MockAuthBackend) GetLoginUrl() string {
	return "http://bklogin.example.com/plain/"
}

// GetUserInfo 获取用户信息
func (b *MockAuthBackend) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	if token == "EverythingIsPermitted" {
		return &UserInfo{ID: "admin"}, nil
	}
	return nil, errors.New("invalid token")
}

var _ AuthBackend = (*MockAuthBackend)(nil)

// NewMockAuthBackend ...
func NewMockAuthBackend() AuthBackend {
	return &MockAuthBackend{}
}
