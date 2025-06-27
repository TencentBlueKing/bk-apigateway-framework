package cmd

import (
	"context"

	"github.com/spf13/cobra"

	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the server version info.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info(context.Background(), version.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
