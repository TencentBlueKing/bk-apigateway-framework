package account

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging/slog-resty"
)

// BkTicketAuthBackend 用于上云版本的用户登录 & 信息获取
type BkTicketAuthBackend struct{}

// GetLoginUrl 获取登录地址
func (b *BkTicketAuthBackend) GetLoginUrl() string {
	return fmt.Sprintf("%s/plain/", config.G.Platform.BkPlatUrl.BkLogin)
}

// GetUserInfo 获取用户信息
func (b *BkTicketAuthBackend) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/user/get_info/", config.G.Platform.BkPlatUrl.BkLogin)

	client := resty.New().SetLogger(slogresty.New(ctx)).SetTimeout(10 * time.Second)

	respData := map[string]any{}
	_, err := client.R().
		SetQueryParams(map[string]string{"bk_ticket": token}).
		ForceContentType("application/json").
		SetResult(&respData).
		Get(url)
	if err != nil {
		return nil, err
	}

	if retCode, cErr := cast.ToIntE(respData["ret"]); cErr != nil {
		return nil, errors.Errorf("get user info api %s return code isn't integer", url)
	} else if retCode != 0 {
		return nil, errors.Errorf("failed to get user info from %s, message: %s", url, respData["msg"])
	}

	data, ok := respData["data"].(map[string]any)
	if !ok {
		return nil, errors.Errorf("failed to get user info from %s, response data not json format", url)
	}
	return &UserInfo{ID: data["username"].(string)}, nil
}

var _ AuthBackend = (*BkTicketAuthBackend)(nil)

// NewBkTicketAuthBackend ...
func NewBkTicketAuthBackend() AuthBackend {
	return &BkTicketAuthBackend{}
}
