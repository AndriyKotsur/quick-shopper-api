package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/AndriyKotsur/quick-booking-api/config"
	"github.com/AndriyKotsur/quick-booking-api/database"
	"github.com/AndriyKotsur/quick-booking-api/internal/server"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migration",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()

		dbConn, err := server.InitDatabase()
		if err != nil {
			log.Fatal().Err(err).Msg("Error: database initialization failed")
		}
		defer dbConn.Close()

		migrator := database.Migrator(dbConn)
		migrator.Up()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
