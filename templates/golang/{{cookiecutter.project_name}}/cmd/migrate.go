package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	// load migration package to register migrations
	_ "bk.tencent.com/{{cookiecutter.project_name}}/pkg/migration"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	database2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/database"
	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/version"
)

// NewMigrateCmd ...
func NewMigrateCmd() *cobra.Command {
	var cfgFile string
	var migrationID string

	migrateCmd := cobra.Command{
		Use:   "migrate",
		Short: "Apply migrations to the database tables.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			if cfg.Platform.Addons.Mysql == nil {
				log.Fatal("mysql config not found, skip migrate...")
			}

			database2.InitDBClient(ctx, cfg.Platform.Addons.Mysql, slog.Default())

			if err = database2.RunMigrate(ctx, migrationID); err != nil {
				log.Fatalf("failed to run migrate: %s", err)
			}
			dbVersion, err := database2.Version(ctx)
			if err != nil {
				log.Fatalf("failed to get database version: %s", err)
			}
			log.Infof(ctx, "migrate success %s\nDatabaseVersion: %s", version.Version(), dbVersion)
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	migrateCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")
	migrateCmd.Flags().StringVar(&migrationID, "migration", "", "migration to apply, blank means latest version")

	return &migrateCmd
}

func init() {
	rootCmd.AddCommand(NewMigrateCmd())
}
