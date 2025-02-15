package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/AndriyKotsur/quick-shopper-api/config"
	"github.com/AndriyKotsur/quick-shopper-api/logger"
)

func InitRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	router.Mount("/v1", v1Router)

	return router
}

func Run() {
	config.LoadConfig()
	logger.InitLogger()
	router := InitRouter()

	server := &http.Server{
		Addr:         ":" + viper.GetString("api.port"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Msgf("Starting server on port: %s", viper.GetString("api.port"))

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Info().Msgf("Receive terminate, graceful shutdown %s", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	server.Shutdown(timeoutContext)
}
