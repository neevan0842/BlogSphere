package main

import (
	"context"

	"github.com/neevan0842/BlogSphere/backend/config"
	"github.com/neevan0842/BlogSphere/backend/database"
	"github.com/neevan0842/BlogSphere/backend/internal"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database connection pool
	pool, err := database.NewDBPool(ctx)
	if err != nil {
		logger.Fatal("Unable to create connection pool: ", err)
	}
	defer pool.Close()

	// Ping the database to ensure connection is established
	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("Unable to connect to database: ", err)
	}
	logger.Info("connected to database pool successfully")

	// API Server
	api := internal.NewAPIServer(config.Envs.ADDR, pool, logger)

	// Run the server
	logger.Fatal(api.Run(api.Mount()))
}
