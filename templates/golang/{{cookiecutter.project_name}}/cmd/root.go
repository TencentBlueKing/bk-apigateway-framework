package cmd

import (
	"context"

	"github.com/spf13/cobra"

	log "github.com/TencentBlueKing/blueapps-go/pkg/logging"
)

var rootCmd = &cobra.Command{
	Use:   "{{cookiecutter.project_name}}",
	Short: "{{cookiecutter.project_name}}",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info(context.Background(), "welcome to use {{cookiecutter.project_name}}, use `{{cookiecutter.project_name}} -h` for help")
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
