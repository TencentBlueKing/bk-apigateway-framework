package cmd

import (
	"context"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/gen"
	"github.com/spf13/cobra"

	config2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/envx"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/filex"
)

// NewGenDefinitionYamlCmd ...
func NewSyncApigatewayCmd() *cobra.Command {
	var cfgFile string
	var docsDir string
	migrateCmd := cobra.Command{
		Use:   "sync_apigateway",
		Short: "sync apigateway",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config2.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}
			baseDir := filex.GetParentDir(docsDir)
			gatewayName := cfg.Platform.AppID
			apiConfig := config2.GetApiConfig(cfg)
			gen.SyncGinGateway(baseDir, gatewayName, apiConfig, true)
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	migrateCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")
	migrateCmd.Flags().StringVar(&docsDir, "docs", envx.Get("DOC_FILE_BASE_DIR", "../"),
		"swagger json docs dir")
	return &migrateCmd
}

func init() {
	rootCmd.AddCommand(NewSyncApigatewayCmd())
}
