package main

import (
	"context"
)

func main() {
	ctx := context.Background()

	cfg := SentinelConfig{
		neo4jConfig: &Neo4jConfig{
			uri:      "bolt://localhost:7687",
			username: "neo4j",
			password: "sentinel",
			database: "neo4j",
		},
		exportPath:   "./export/",
		defaultLimit: 5,
	}

	sentinel, err := NewSentinel(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer sentinel.stop(ctx)

	sentinel.logger.Info("Purging DB...")
	if err := sentinel.dbConn.PurgeDB(ctx); err != nil {
		panic(err)
	}

	sentinel.logger.Info("Mapping folder...")
	if err = sentinel.mapper.MapFolder(ctx, "examples"); err != nil {
		panic(err)
	}

	fullProjection := Projection{
		graphName:      DefaultGraphProjectionName,
		nodeProjection: DefaultNodeProjectionNames,
		relProjection:  DefaultRelProjectionNames,
	}

	sentinel.logger.Info("Projecting graph...")
	sentinel.analyzer.ProjectGraph(ctx, fullProjection)

	sentinel.logger.Info("Analyzing centrality...")

	betweenesList := sentinel.analyzer.Centrality(ctx, "betweenness", fullProjection, 5)
	closenessList := sentinel.analyzer.Centrality(ctx, "closeness", fullProjection, 5)
	eigenvectorList := sentinel.analyzer.Centrality(ctx, "eigenvector", fullProjection, 5)
	degreeList := sentinel.analyzer.Centrality(ctx, "degree", fullProjection, 5)

	sentinel.logger.Info("Exporting...")

	// sentinel.exporter.logNodes(betweenesList)
	// sentinel.exporter.logNodes(closenessList)
	// sentinel.exporter.logNodes(eigenvectorList)
	// sentinel.exporter.logNodes(degreeList)

	sentinel.exporter.writeToFile(betweenesList, "betweenness.json")
	sentinel.exporter.writeToFile(closenessList, "closeness.json")
	sentinel.exporter.writeToFile(eigenvectorList, "eigenvector.json")
	sentinel.exporter.writeToFile(degreeList, "degree.json")

	// TODO: community: CALL gds.louvain.stream('fullGraph', {})

}
