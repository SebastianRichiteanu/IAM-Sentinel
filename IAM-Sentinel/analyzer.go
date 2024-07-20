package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"
)

type (
	Analyzer struct {
		logger       *logrus.Logger
		dbConn       *Neo4jConn
		defaultLimit int
	}

	Projection struct {
		graphName      string
		nodeProjection []string
		relProjection  []string
	}
)

func NewAnalyzer(logger *logrus.Logger, db *Neo4jConn, defaultLimit int) (*Analyzer, error) {
	a := Analyzer{
		logger:       logger,
		dbConn:       db,
		defaultLimit: defaultLimit,
	}
	return &a, nil
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
	RETURN gds.util.asNode(nodeId) as node, score as score
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
		nodeField, ok := result.Get("node")
		if !ok {
			a.logger.WithField("centralityType", centralityType).Warn("could not get node result")
			continue
		}

		node, ok := nodeField.(neo4j.Node)
		if !ok {
			a.logger.WithField("centralityType", centralityType).Warn("could not cast node result")
			continue
		}

		score, ok := result.Get("score")
		if !ok {
			a.logger.WithField("centralityType", centralityType).Warn("could not get score")
			continue
		}

		newNode := CentralityNode{
			Score: score.(float64),
			Node:  &node,
		}

		nodes = append(nodes, newNode)
	}

	return nodes
}

func (a *Analyzer) Community(ctx context.Context, communityType string, projection Projection) CommunityMap {
	query := fmt.Sprintf(`
	CALL gds.%s.stream($graph_name, {}) 
	YIELD nodeId, communityId
	RETURN gds.util.asNode(nodeId) as node, communityId as communityId
	`, communityType)

	params := map[string]any{
		"graph_name": projection.graphName,
	}

	results, err := a.dbConn.ExecuteQueryRead(ctx, query, params)
	if err != nil {
		a.logger.Error(err)
	}

	communities := make(CommunityMap)

	for _, result := range results {
		nodeField, ok := result.Get("node")
		if !ok {
			a.logger.WithField("communityType", communityType).Warn("could not get node result")
			continue
		}

		node, ok := nodeField.(neo4j.Node)
		if !ok {
			a.logger.WithField("communityType", communityType).Warn("could not cast node result")
			continue
		}

		communityIdField, ok := result.Get("communityId")
		if !ok {
			a.logger.WithField("communityType", communityType).Warn("could not get communityID")
			continue
		}
		communityId, ok := communityIdField.(int64)
		if !ok {
			a.logger.WithField("communityType", communityType).Warn("could not cast communityID")
			continue
		}

		communities[communityId] = append(communities[communityId], &node)
	}

	return communities
}
