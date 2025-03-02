package cmd

import (
	"github.com/spf13/cobra"

	"github.com/AndriyKotsur/quick-shopper-api/database"
	"github.com/AndriyKotsur/quick-shopper-api/internal/server"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migration",
	Run: func(cmd *cobra.Command, args []string) {
		db := server.InitDatabase()

		migrator := database.Migrator(db)
		migrator.Up()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
