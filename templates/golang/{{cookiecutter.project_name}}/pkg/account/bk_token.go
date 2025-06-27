package account

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging/slog-resty"
)

// BkTokenAuthBackend 用于社区开源版本的用户登录 & 信息获取
type BkTokenAuthBackend struct{}

// GetLoginUrl 获取登录地址
func (b BkTokenAuthBackend) GetLoginUrl() string {
	return fmt.Sprintf("%s/plain/", config.G.Platform.BkPlatUrl.BkLogin)
}

// GetUserInfo 获取用户信息
func (b BkTokenAuthBackend) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/accounts/get_user/", config.G.Platform.BkPlatUrl.BkLogin)

	client := resty.New().SetLogger(slogresty.New(ctx)).SetTimeout(10 * time.Second)

	respData := map[string]any{}
	_, err := client.R().
		SetQueryParams(map[string]string{"bk_token": token}).
		ForceContentType("application/json").
		SetResult(&respData).
		Get(url)
	if err != nil {
		return nil, err
	}

	if retCode, cErr := cast.ToIntE(respData["code"]); cErr != nil {
		return nil, errors.Errorf("get user info api %s return code isn't integer", url)
	} else if retCode != 0 {
		return nil, errors.Errorf("failed to get user info from %s, message: %s", url, respData["message"])
	}

	data, ok := respData["data"].(map[string]any)
	if !ok {
		return nil, errors.Errorf("failed to get user info from %s, response data not json format", url)
	}
	return &UserInfo{ID: data["username"].(string)}, nil
}

var _ AuthBackend = (*BkTokenAuthBackend)(nil)

// NewBkTokenAuthBackend ...
func NewBkTokenAuthBackend() AuthBackend {
	return &BkTokenAuthBackend{}
}
