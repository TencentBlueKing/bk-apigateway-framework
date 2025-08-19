package cmd

import (
	"context"

	"github.com/spf13/cobra"

	log "github.com/TencentBlueKing/blueapps-go/pkg/logging"
)

var rootCmd = &cobra.Command{
	Use:   "code-gw-go-demo",
	Short: "Golang Gin Template V2",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info(context.Background(), "welcome to use code-gw-go-demo, use `code-gw-go-demo -h` for help")
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
