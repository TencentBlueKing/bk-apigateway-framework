package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/otel"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/router"
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

			// 初始化 Logger
			if err = initLogger(&cfg.Service.Log); err != nil {
				logging.Fatalf("failed to init logging: %s", err)
			}
			// 初始化增强服务客户端
			if err = initAddons(ctx, cfg); err != nil {
				logging.Fatalf("failed to init addons: %s", err)
			}
			// 初始化 OpenTelemetry
			if cfg.Platform.Addons.BkOtel != nil {
				shutdown, sErr := otel.InitTracer(ctx, cfg.Platform.Addons.BkOtel, otel.GenServiceName("web"))
				if sErr != nil {
					logging.Fatalf("failed to init OpenTelemetry: %s", sErr)
				}
				defer func() {
					if err = shutdown(ctx); err != nil {
						logging.Fatalf("failed to shutdown OpenTelemetry: %s", err)
					}
				}()
			}

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

func init() {
	rootCmd.AddCommand(NewWebServerCmd())
}
