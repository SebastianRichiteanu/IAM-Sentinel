package main

import (
	"context"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultGraphProjectionName = "fullGraph"
)

var (
	DefaultNodeProjectionNames = []string{"IAM"}
	DefaultRelProjectionNames  = []string{"*"}
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

	// Clear
	if err := db.PurgeDB(ctx); err != nil {
		panic(err)
	}

	// Parse
	parser, err := NewParser(logger)
	if err != nil {
		panic(err)
	}

	mapper, err := NewResourceMapper(logger, db, parser)
	if err != nil {
		panic(err)
	}

	if err = mapper.MapFolder(ctx, "examples"); err != nil {
		panic(err)
	}

	// Analyze
	analyzer, err := NewAnalyzer(logger, db)
	if err != nil {
		panic(err)
	}

	fullProjection := Projection{
		graphName:      DefaultGraphProjectionName,
		nodeProjection: DefaultNodeProjectionNames,
		relProjection:  DefaultRelProjectionNames,
	}

	analyzer.ProjectGraph(ctx, fullProjection)

	betweenesList := analyzer.Centrality(ctx, "betweenness", fullProjection, 5)
	closenessList := analyzer.Centrality(ctx, "closeness", fullProjection, 5)
	eigenvectorList := analyzer.Centrality(ctx, "eigenvector", fullProjection, 5)
	degreeList := analyzer.Centrality(ctx, "degree", fullProjection, 5)

	exportNodes(logger, betweenesList)
	exportNodes(logger, closenessList)
	exportNodes(logger, eigenvectorList)
	exportNodes(logger, degreeList)
}

// TODO: make exporter component?
func exportNodes(logger *log.Logger, nodes CentralityNodes) {
	for _, node := range nodes {
		jsonBytes, err := json.Marshal(node)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info(string(jsonBytes))
	}
}
