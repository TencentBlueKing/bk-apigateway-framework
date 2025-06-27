package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_contrib/gen"
	"github.com/spf13/cobra"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/config"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/router"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/envx"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/filex"
)

// NewGenDefinitionYamlCmd ...
func NewGenResourceYamlCmd() *cobra.Command {
	var cfgFile string
	var docsDir string
	migrateCmd := cobra.Command{
		Use:   "gen_resources_yaml",
		Short: "generate definition.yaml ",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			_, err := config.Load(ctx, cfgFile)
			if err != nil {
				logging.Fatalf("failed to load config: %s", err)
			}
			engine := router.New(logging.GetLogger("gin"))
			docPath := docsDir + "/swagger.json"
			yaml := gen.GenResourceYamlFromSwaggerJson(docPath, engine)
			logging.Infof(ctx, "gen resource yaml success:\n %s", yaml)
			resourcesFilePath := filepath.Join(filex.GetParentDir(docsDir), "resources.yaml")
			// 生成资源配置yaml文件
			err = os.WriteFile(resourcesFilePath, []byte(yaml), 0o644)
			if err != nil {
				logging.Fatalf("failed to write file: %s", err)
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
	rootCmd.AddCommand(NewGenResourceYamlCmd())
}
