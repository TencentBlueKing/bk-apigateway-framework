// Package account 提供不同版本的用户认证后端
package account

import (
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common"
)

// GetAuthBackend 通过 Token 类型获取 AuthBackend
func GetAuthBackend() AuthBackend {
	switch common.UserTokenKey {
	case "bk_ticket":
		return NewBkTicketAuthBackend()
	case "bk_token":
		return NewBkTokenAuthBackend()
	default:
		// 默认为 BkToken
		return NewBkTokenAuthBackend()
	}
}
