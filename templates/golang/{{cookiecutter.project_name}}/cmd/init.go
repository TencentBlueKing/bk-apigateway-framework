package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/cache/memory"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/infras/database"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/infras/redis"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
)

func initLogger(cfg *config.LogConfig) error {
	// 自动创建日志目录
	if err := os.MkdirAll(cfg.Dir, os.ModePerm); err != nil {
		// 只有当错误不是 “目录已存在” 时，需要抛出错误
		if !os.IsExist(err) {
			return errors.Wrapf(err, "creating log dir %s", cfg.Dir)
		}
	}

	// 输出位置
	writerName := "file"
	if cfg.ForceToStdout {
		writerName = "stdout"
	}

	// 初始化默认 Logger
	loggerName := "default"
	if err := logging.InitLogger(loggerName, &logging.Options{
		Level: cfg.Level,
		// 输出到标准输出时 Text 会更加友好
		HandlerName:  lo.Ternary(writerName == "stdout", "text", "json"),
		WriterName:   writerName,
		WriterConfig: map[string]string{"filename": filepath.Join(cfg.Dir, loggerName+".log")},
	}); err != nil {
		return errors.Wrapf(err, "creating logger %s", loggerName)
	}

	// 初始化 Gorm Logger
	loggerName = "gorm"
	if err := logging.InitLogger(loggerName, &logging.Options{
		Level:        logging.GormLogLevel,
		HandlerName:  "json",
		WriterName:   "file",
		WriterConfig: map[string]string{"filename": filepath.Join(cfg.Dir, loggerName+".log")},
	}); err != nil {
		return errors.Wrapf(err, "creating logger %s", loggerName)
	}

	// 初始化 Gin Logger
	loggerName = "gin"
	if err := logging.InitLogger(loggerName, &logging.Options{
		Level:        logging.GinLogLevel,
		HandlerName:  "json",
		WriterName:   "file",
		WriterConfig: map[string]string{"filename": filepath.Join(cfg.Dir, loggerName+".log")},
	}); err != nil {
		return errors.Wrapf(err, "creating logger %s", loggerName)
	}

	return nil
}

// 根据增强服务配置，初始化各类客户端
func initAddons(ctx context.Context, cfg *config.Config) error {
	// 初始化 DB Client
	database.InitDBClient(ctx, cfg.Platform.Addons.Mysql, logging.GetLogger("gorm"))

	// 初始化 Redis Client
	if cfg.Platform.Addons.Redis != nil {
		redis.InitRedisClient(ctx, cfg.Platform.Addons.Redis)
	}

	// 初始化缓存
	memory.InitCache(cfg.Service.MemoryCacheSize)

	return nil
}
