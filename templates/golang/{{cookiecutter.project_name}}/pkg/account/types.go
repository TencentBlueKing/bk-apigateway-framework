package account

import "context"

// UserInfo 用户信息
type UserInfo struct {
	ID string
}

// AuthBackend 认证后端
type AuthBackend interface {
	// GetLoginUrl 获取登录地址
	GetLoginUrl() string
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)
}
