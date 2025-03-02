package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/AndriyKotsur/quick-shopper-api/config"
	"github.com/AndriyKotsur/quick-shopper-api/database"
	"github.com/AndriyKotsur/quick-shopper-api/internal/server"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run seeding",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()

		dbConn, err := server.InitDatabase()
		if err != nil {
			log.Fatal().Err(err).Msg("Error: database initialization failed")
		}
		defer dbConn.Close()

		seeder := database.Seeder(dbConn)
		seeder.SeedUsers()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
