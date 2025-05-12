// Package objstorage 提供对象存储相关封装，目前接入的是蓝盾制品库（bkrepo）
// 如果 SaaS 开发者需要使用其他云对象存储（如 COS，S3, Ceph 等），可参考相关实现
package objstorage

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	config2 "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/infras/otel/otel-resty"
	log "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging/slog-resty"
)

// Client 获取对象存储客户端
func NewClient(ctx context.Context) *BkGenericRepoClient {
	if !IsBkRepoAvailable() {
		log.Error(ctx, "bkrepo is not available")
		return nil
	}

	cli := newBkGenericRepoClient(config2.G.Platform.Addons.BkRepo)
	// 根据 ctx 设置 Logger，以支持记录 Request / Span / Trace ID 等信息
	cli.client = cli.client.SetLogger(slogresty.New(ctx))
	return cli
}

// IsBkRepoAvailable 判断蓝盾制品仓库是否可用
func IsBkRepoAvailable() bool {
	return config2.G.Platform.Addons.BkRepo != nil
}

// 初始化客户端 ...
func newBkGenericRepoClient(cfg *config2.BkRepoConfig) *BkGenericRepoClient {
	// 使用连接池
	transport := &http.Transport{
		MaxIdleConns: 5, MaxIdleConnsPerHost: 5, IdleConnTimeout: 30 * time.Second,
	}

	client := resty.New().
		SetTransport(transport).
		SetBaseURL(strings.TrimSuffix(cfg.EndpointUrl, "/")).
		SetTimeout(60*time.Second).
		SetBasicAuth(cfg.Username, cfg.Password).
		SetRetryCount(2).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
			AddRetryCondition(
				func(response *resty.Response, err error) bool {
					// Retry on 5xx status codes
					return response.StatusCode() >= http.StatusInternalServerError
				},
			).
		// OpenTelemetry 相关中间件
		OnBeforeRequest(otelresty.RequestMiddleware).
		OnAfterResponse(otelresty.ResponseMiddleware)

	return &BkGenericRepoClient{cfg: cfg, client: client}
}
