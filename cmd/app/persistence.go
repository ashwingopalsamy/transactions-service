package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func InitDB(c *EnvCfg) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check if the database is reachable
	if err := healthCheckDB(dbPool); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("database is not reachable: %w", err)
	}

	log.Info().Msg("database connection established")
	return dbPool, nil
}

func healthCheckDB(dbPool *pgxpool.Pool) error {
	if err := dbPool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	log.Info().Msg("database is reachable")
	return nil
}
