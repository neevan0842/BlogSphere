package main

import (
	"context"

	"github.com/neevan0842/BlogSphere/backend/config"
	"github.com/neevan0842/BlogSphere/backend/database"
	"github.com/neevan0842/BlogSphere/backend/internal"
	"github.com/neevan0842/BlogSphere/backend/logger"
)

func main() {
	ctx := context.Background()

	// Logger
	logger.Init()
	defer logger.Sync()
	log := logger.Get()

	// Database connection pool
	pool, err := database.NewDBPool(ctx)
	if err != nil {
		log.Fatal("Unable to create connection pool: ", err)
	}
	defer pool.Close()

	// Ping the database to ensure connection is established
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	log.Info("connected to database pool successfully")

	// API Server
	server := internal.NewAPIServer(config.Envs.ADDR, pool, log)

	// Run the server
	log.Fatal(server.Run(server.Mount()))
}
