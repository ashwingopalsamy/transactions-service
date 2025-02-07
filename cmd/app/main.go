package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultWebPort = 8080

type EnvCfg struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     int
	Port       int
}

func main() {
	// Initialize logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load environmental config
	cfg := loadEnvConfig()

	// Initialize database
	dbPool, err := InitDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer dbPool.Close()

	// Setup server
	router := NewRouter()
	server := NewServer(router, withPort(cfg.Port))
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start web server")
		}
	}()

	// Setup graceful shutdown
	gracefulShutdown(server)
}

func loadEnvConfig() *EnvCfg {
	return &EnvCfg{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		Port:       getEnvAsInt("PORT", defaultWebPort),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal().Msgf("invalid port value: %s", value)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}

func gracefulShutdown(server HTTPServer) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan

	log.Info().Msg("Received termination signal, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shut down server gracefully")
	}
	log.Info().Msg("Server gracefully stopped")
}
