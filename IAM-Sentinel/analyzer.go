package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"
)

type (
	Analyzer struct {
		logger *logrus.Logger
		dbConn *Neo4jConn
	}

	Projection struct {
		graphName      string
		nodeProjection []string
		relProjection  []string
	}
)

func NewAnalyzer(logger *logrus.Logger, db *Neo4jConn) (*Analyzer, error) {
	a := Analyzer{
		logger: logger,
		dbConn: db,
	}
	return &a, nil
}

func (a *Analyzer) GetNodeList(ctx context.Context, ids []int64) CentralityNodes {
	var nodes CentralityNodes

	query := `MATCH(n) WHERE id(n) in $ids RETURN n`

	params := map[string]any{"ids": ids}

	results, err := a.dbConn.ExecuteQueryRead(ctx, query, params)
	if err != nil {
		a.logger.Error(err)
	}

	for _, result := range results {
		nodeAsResult, ok := result.Get("n")
		if !ok {
			a.logger.Error("could not get node for node map: %w", err)
			continue
		}

		node, ok := nodeAsResult.(neo4j.Node)
		if !ok {
			a.logger.Error("could not map result to node: %w", err)
			continue
		}

		jsonNode := NodeJSON{
			ID:     node.GetId(),
			Labels: node.Labels,
			Props:  node.Props,
		}

		nodes = append(nodes, jsonNode)
	}
	return nodes
}

func (a *Analyzer) ProjectGraph(ctx context.Context, projection Projection) {
	query := `CALL gds.graph.project($graph_name,$node_projection,$rel_projection)`

	params := map[string]any{
		"graph_name":      projection.graphName,
		"node_projection": projection.nodeProjection,
		"rel_projection":  projection.relProjection,
	}

	if _, err := a.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
		a.logger.Error(err)
	}
}

func (a *Analyzer) Centrality(ctx context.Context, centralityType string, projection Projection, noNodes int) CentralityNodes {
	query := fmt.Sprintf(`
	CALL gds.%s.stream($graph_name, {}) YIELD nodeId, score
	RETURN nodeId, score
	ORDER BY score DESC
	LIMIT $limit_nr`, centralityType)

	params := map[string]any{
		"graph_name": projection.graphName,
		"limit_nr":   noNodes,
	}

	results, err := a.dbConn.ExecuteQueryRead(ctx, query, params)
	if err != nil {
		a.logger.Error(err)
	}

	var ids []int64

	for _, result := range results {
		ids = append(ids, result.Values[0].(int64))
	}

	return a.GetNodeList(ctx, ids)
}
