package main

import (
	"github.com/rs/zerolog/log"

	"github.com/AndriyKotsur/quick-shopper-api/cmd"
)

func main() {
	log.Info().Msg("Starting application...")
	cmd.Execute()
}
