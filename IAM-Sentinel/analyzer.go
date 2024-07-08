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

func (a *Analyzer) GetNodeJSON(ctx context.Context, id int64) *NodeJSON {
	query := `MATCH(n) WHERE id(n) = $id RETURN n`
	params := map[string]any{"id": id}

	results, err := a.dbConn.ExecuteQueryRead(ctx, query, params)
	if err != nil {
		a.logger.Error(err)
	}

	if len(results) == 0 {
		a.logger.Error("could not find node", "id", id)
		return nil
	}

	nodeAsResult, ok := results[0].Get("n")
	if !ok {
		a.logger.Error("could not get node for node map: %w", err)
		return nil
	}

	node, ok := nodeAsResult.(neo4j.Node)
	if !ok {
		a.logger.Error("could not map result to node: %w", err)
		return nil
	}

	return &NodeJSON{
		ID:     node.GetId(),
		Labels: node.Labels,
		Props:  node.Props,
	}
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

func (a *Analyzer) Centrality(ctx context.Context, centralityType string, projection Projection, noNodes int) []CentralityNode {
	query := fmt.Sprintf(`
	CALL gds.%s.stream($graph_name, {}) YIELD nodeId, score
	RETURN nodeId as id, score as score
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

	var nodes []CentralityNode

	for _, result := range results {
		id, ok := result.Get("id")
		if !ok {
			continue
		}
		score, ok := result.Get("score")
		if !ok {
			continue
		}

		node := a.GetNodeJSON(ctx, id.(int64))
		if node == nil {
			continue
		}

		newNode := CentralityNode{
			Score: score.(float64),
			Node:  *node,
		}

		nodes = append(nodes, newNode)
	}

	return nodes
}
