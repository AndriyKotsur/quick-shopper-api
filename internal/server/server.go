package server

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/AndriyKotsur/quick-shopper-api/config"
	"github.com/AndriyKotsur/quick-shopper-api/internal/database"
	"github.com/AndriyKotsur/quick-shopper-api/internal/domain/authentication"
	"github.com/AndriyKotsur/quick-shopper-api/internal/domain/user"
	logger "github.com/AndriyKotsur/quick-shopper-api/logger"
	db "github.com/AndriyKotsur/quick-shopper-api/third_party/database"
)

func InitDatabase() (*sql.DB, error) {
	dbConfig := db.Config{
		Host:     viper.GetString("db.host"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Name:     viper.GetString("db.name"),
		Port:     viper.GetString("db.port"),
	}

	database, err := db.Connect(dbConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	return database, nil
}
func InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger)

	return router
}

func Run() {
	config.LoadConfig()
	logger.InitLogger()

	dbConn, err := InitDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Database initialization failed")
	}
	defer dbConn.Close()

	router := InitRouter()
	dbQueries := database.New(dbConn)

	v1Router := chi.NewRouter()
	authentication.RegisterAuthEndpoints(v1Router, dbQueries)
	user.RegisterUserEndpoints(v1Router, dbQueries)

	router.Mount("/v1", v1Router)

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan
	log.Info().Msgf("Receive terminate, graceful shutdown %s", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	server.Shutdown(timeoutContext)
}
