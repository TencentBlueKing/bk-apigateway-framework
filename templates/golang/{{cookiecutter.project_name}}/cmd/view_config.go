package cmd

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/config"
	log "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
)

// NewViewConfigCmd ...
func NewViewConfigCmd() *cobra.Command {
	var cfgFile string
	var verbose bool

	viewCfgCmd := cobra.Command{
		Use:   "view-config",
		Short: "View service configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			// 加载配置
			cfg, err := config.Load(context.Background(), cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			if verbose {
				spew.Dump(cfg)
				return
			}

			data, _ := yaml.Marshal(cfg)
			fmt.Println(string(data))
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	viewCfgCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")
	viewCfgCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show more details")

	return &viewCfgCmd
}

func init() {
	rootCmd.AddCommand(NewViewConfigCmd())
}
