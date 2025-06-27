package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/async"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/otel"
	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
)

// NewSchedulerCmd 用于创建定时任务调度器启动命令
// 需要注意的是：为避免重复执行定时任务，需要确保同时只有一个 scheduler 正在运行
// 如果希望同时启动多个 scheduler 启动，则需要添加诸如 redis / zk 这样的分布式锁
func NewSchedulerCmd() *cobra.Command {
	var cfgFile string

	schedulerCmd := cobra.Command{
		Use:   "scheduler",
		Short: "Execute tasks based on cron expressions, please ensure only one running scheduler.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			// 初始化 Logger
			if err = initLogger(&cfg.Service.Log); err != nil {
				log.Fatalf("failed to init logging: %s", err)
			}
			// 初始化增强服务客户端
			if err = initAddons(ctx, cfg); err != nil {
				log.Fatalf("failed to init addons: %s", err)
			}
			// 初始化 OpenTelemetry
			if cfg.Platform.Addons.BkOtel != nil {
				shutdown, sErr := otel.InitTracer(ctx, cfg.Platform.Addons.BkOtel, otel.GenServiceName("scheduler"))
				if sErr != nil {
					log.Fatalf("failed to init OpenTelemetry: %s", sErr)
				}
				defer func() {
					if err = shutdown(ctx); err != nil {
						log.Fatalf("failed to shutdown OpenTelemetry: %s", err)
					}
				}()
			}

			// 初始化 task server
			async.InitTaskScheduler(ctx)

			srv := async.Scheduler()
			// 加载周期任务
			if err = srv.LoadTasks(); err != nil {
				log.Fatal(err.Error())
			}
			// 启用调度服务器
			srv.Run()
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	schedulerCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")

	return &schedulerCmd
}

func init() {
	rootCmd.AddCommand(NewSchedulerCmd())
}
