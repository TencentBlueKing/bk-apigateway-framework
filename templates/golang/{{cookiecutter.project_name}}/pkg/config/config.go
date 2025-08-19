// Package config 管理蓝鲸 SaaS 配置项，支持从配置文件 / 环境变量中读取配置
package config

import (
	"context"

	"github.com/TencentBlueKing/blueapps-go/pkg/config"
	"github.com/pkg/errors"
)

var G *SvcConfig

// Load 加载配置
func Load(ctx context.Context, cfgFile string) (*SvcConfig, error) {
	cfg, err := config.Load(ctx, cfgFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "load config file %s failed", cfgFile)
	}
	G = &SvcConfig{
		Config: cfg,
	}
	return G, nil
}
