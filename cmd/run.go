package cmd

import (
	"github.com/spf13/cobra"

	"github.com/AndriyKotsur/quick-shopper-api/internal/server"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}
