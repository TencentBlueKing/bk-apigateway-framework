package cmd

import (
	"context"

	"github.com/spf13/cobra"

	log "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
)

var rootCmd = &cobra.Command{
	Use:   "{{cookiecutter.project_name}}",
	Short: "Golang Gin Template V2",
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
