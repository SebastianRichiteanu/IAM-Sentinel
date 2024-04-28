package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func initLogger() *log.Logger {
	logger := log.New()
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	return logger
}

func main() {
	ctx := context.Background()
	logger := initLogger()

	logger.Info("Starting...")

	db, err := InitializeDB(ctx, Neo4jConfig{
		uri:      "bolt://localhost:7687",
		username: "neo4j",
		password: "sentinel",
		database: "neo4j"}) // TODO: move this to env params or config or something else :)
	if err != nil {
		panic(err)
	}
	defer db.Stop(ctx)

	if err := db.PurgeDB(ctx); err != nil {
		panic(err)
	}

	parser, err := NewParser(logger)
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile("example.json")
	if err != nil {
		panic(err)
	}

	mapper, err := NewResourceMapper(logger, db, parser)
	if err != nil {
		panic(err)
	}

	if err = mapper.MapFile(ctx, data); err != nil {
		panic(err)
	}
}
