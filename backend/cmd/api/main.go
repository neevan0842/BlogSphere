package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/neevan0842/BlogSphere/backend/config"
	"github.com/neevan0842/BlogSphere/backend/internal"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	conn, err := pgx.Connect(ctx, config.Envs.DATABASE_URL)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("connected to database successfully")
	defer conn.Close(ctx)

	// API Server
	api := internal.NewAPIServer(config.Envs.ADDR, conn, logger)

	// Run the server
	logger.Fatal(api.Run(api.Mount()))
}
