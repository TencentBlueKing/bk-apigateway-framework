package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/TencentBlueKing/blueapps-go/pkg/infras/database"
	"github.com/TencentBlueKing/blueapps-go/pkg/infras/redis"
	"github.com/TencentBlueKing/blueapps-go/pkg/logging"
	"github.com/spf13/cobra"

	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/config"
	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/router"
)

// NewWebServerCmd ...
func NewWebServerCmd() *cobra.Command {
	var cfgFile string

	wsCmd := cobra.Command{
		Use:   "webserver",
		Short: "Start the HTTP server.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config.Load(ctx, cfgFile)
			if err != nil {
				logging.Fatalf("failed to load config: %s", err)
			}
			// 初始化增强服务客户端
			initAddons(ctx, cfg)

			// 启动 Web 服务
			logging.Infof(ctx, "Starting server at http://0.0.0.0:%d", config.G.Service.Server.Port)
			srv := &http.Server{
				Addr:    ":" + strconv.Itoa(cfg.Service.Server.Port),
				Handler: router.New(logging.GetLogger("gin")),
			}
			go func() {
				if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logging.Fatalf("Start server failed: %s", err)
				}
			}()

			// 等待中断信号以优雅地关闭服务器
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit

			srvCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Service.Server.GraceTimeout)*time.Second)
			defer cancel()

			logging.Info(ctx, "Shutdown server ...")
			if err = srv.Shutdown(srvCtx); err != nil {
				logging.Fatalf("Shutdown server failed: %s", err)
			}
			logging.Info(ctx, "Server exiting")
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	wsCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")

	return &wsCmd
}

// 根据增强服务配置，初始化各类客户端
// 注意：/healthz 会探测已配置的增强服务，因此配置了的增强服务必须在此初始化客户端
func initAddons(ctx context.Context, cfg *config.SvcConfig) {
	// 初始化 DB Client
	if cfg.Platform.Addons.Mysql != nil {
		database.InitDBClient(ctx, cfg.Platform.Addons.Mysql, logging.GetLogger("gorm"))
	}

	// 初始化 Redis Client
	if cfg.Platform.Addons.Redis != nil {
		redis.InitRedisClient(ctx, cfg.Platform.Addons.Redis)
	}
}

func init() {
	rootCmd.AddCommand(NewWebServerCmd())
}
