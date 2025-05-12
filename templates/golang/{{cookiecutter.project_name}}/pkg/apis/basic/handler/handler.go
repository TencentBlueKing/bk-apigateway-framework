// Package handler ...
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/apis/basic/serializer"
	probe2 "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/common/probe"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/version"
)

// Ping ...
//
//	@Summary	服务探活
//	@Tags		basic
//	@Produce	text/plain
//	@Success	200	{string}	string	pong
//	@Router		/ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Healthz ...
//
//	@Summary	提供服务健康状态
//	@Tags		basic
//	@Param		token	query		string	true	"healthz api token"
//	@Success	200		{object}	serializer.HealthzResponse
//	@Router		/healthz [get]
func Healthz(c *gin.Context) {
	ctx := c.Request.Context()

	healthy, fatal := true, false
	var results []probe2.Result

	probes := []probe2.HealthProbe{
		probe2.NewGin(),
		probe2.NewMysql(),
		probe2.NewRedis(),
		probe2.NewBkRepo(),
	}
	for _, p := range probes {
		ret := p.Perform(ctx)
		if ret == nil {
			continue
		}

		// 任意探针失败，则为不健康
		healthy = healthy && ret.Healthy
		// 任意核心组件探针失败，则为致命异常
		fatal = fatal || (ret.Core && !ret.Healthy)
		results = append(results, *ret)
	}
	respData := serializer.HealthzResponse{
		Time:    time.Now().Format(time.RFC3339),
		Healthy: healthy,
		Fatal:   fatal,
		Results: results,
	}
	// 如果核心服务不可用，应该返回 503 而非 200
	c.JSON(lo.Ternary(fatal, http.StatusServiceUnavailable, http.StatusOK), respData)
}

// Version ...
//
//	@Summary	服务版本信息
//	@Tags		basic
//	@Success	200	{object}	serializer.VersionResponse
//	@Router		/version [get]
func Version(c *gin.Context) {
	respData := serializer.VersionResponse{
		Version:     version.AppVersion,
		GitCommit:   version.GitCommit,
		BuildTime:   version.BuildTime,
		TmplVersion: version.TmplVersion,
		GoVersion:   version.GoVersion,
	}
	c.JSON(http.StatusOK, respData)
}

// Metrics ...
//
//	@Summary	Prometheus 指标
//	@Tags		basic
//	@Produce	text/plain
//	@Param		token	query		string	true	"metrics api token"
//	@Success	200		{string}	string	metrics
//	@Router		/metrics [get]
func Metrics() {} // nolint
