package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/gen"
	"github.com/spf13/cobra"

	config2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/envx"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/filex"
)

// NewGenDefinitionYamlCmd ...
func NewGenDefinitionYamlCmd() *cobra.Command {
	var cfgFile string
	var docsDir string
	migrateCmd := cobra.Command{
		Use:   "gen_definition_yaml",
		Short: "generate definition.yaml ",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config2.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}
			apiConfig := config2.GetApiConfig(cfg)
			yaml := gen.GenDefinitionYaml(apiConfig)
			log.Infof(ctx, "gen definition yaml success:\n %s", yaml)
			definitionFilePath := filepath.Join(filex.GetParentDir(docsDir), "definition.yaml")
			// 生成资源配置yaml文件
			err = os.WriteFile(definitionFilePath, []byte(yaml), 0o644)
			if err != nil {
				log.Fatalf("failed to write file: %s", err)
			}

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
	rootCmd.AddCommand(NewGenDefinitionYamlCmd())
}
