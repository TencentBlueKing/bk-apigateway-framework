package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/gen"
	log "github.com/TencentBlueKing/blueapps-go/pkg/logging"
	"github.com/TencentBlueKing/blueapps-go/pkg/utils/envx"
	"github.com/spf13/cobra"

	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/config"
	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/router"
	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/utils"
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
			cfg, err := config.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}
			apiConfig := config.GetApiConfig(cfg)
			engine := router.New(log.GetLogger("gin"))

			// 校验 docsDir，防止路径遍历导致任意文件写入
			absDocsDir, err := filepath.Abs(docsDir)
			if err != nil {
				log.Fatalf("failed to resolve docs dir: %s", err)
			}
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current working directory: %s", err)
			}
			rel, err := filepath.Rel(cwd, absDocsDir)
			if err != nil {
				log.Fatalf("failed to validate docs dir: %s", err)
			}
			if strings.HasPrefix(rel, "..") {
				log.Fatalf("invalid docs dir %q: path traversal detected", docsDir)
			}

			docPath := docsDir + "/swagger.json"
			yaml := gen.GenDefinitionYaml(apiConfig, docPath, engine)
			log.Infof(ctx, "gen definition yaml success:\n %s", yaml)
			definitionFilePath := filepath.Join(utils.GetParentDir(docsDir), "definition.yaml")
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
	migrateCmd.Flags().StringVar(&docsDir, "docs", envx.Get("DOC_FILE_BASE_DIR", "./docs"),
		"swagger json docs dir")
	return &migrateCmd
}

func init() {
	rootCmd.AddCommand(NewGenDefinitionYamlCmd())
}
