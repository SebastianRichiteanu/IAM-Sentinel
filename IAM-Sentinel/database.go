package main

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	queryPurgeDB        = "MATCH(n) DETACH DELETE(n)"
	createKeyConstraint = "CREATE CONSTRAINT FOR (n) REQUIRE n.key IS UNIQUE;"

	timeout = time.Second * 30
)

type Neo4jConn struct {
	config Neo4jConfig
	driver neo4j.DriverWithContext
}

type Neo4jConfig struct {
	username string
	password string
	uri      string
	database string
}

func InitializeDB(ctx context.Context, config Neo4jConfig) (*Neo4jConn, error) {
	driver, err := neo4j.NewDriverWithContext(config.uri, neo4j.BasicAuth(config.username, config.password, ""))
	if err != nil {
		return nil, fmt.Errorf("could not create new neo4j driver: %w", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create verify neo4j driver connectivity: %w", err)
	}

	return &Neo4jConn{config: config, driver: driver}, nil
}

func (neo Neo4jConn) Stop(ctx context.Context) {
	neo.driver.Close(ctx)
}

func (neo Neo4jConn) ExecuteQueryRead(ctx context.Context, query string, params map[string]any) ([]*neo4j.Record, error) {
	session := neo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: neo.config.database})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, getTransactionFunc(ctx, query, params), neo4j.WithTxTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("could not execute neo4j session read: %w", err)
	}

	return result.([]*neo4j.Record), nil
}

func (neo Neo4jConn) ExecuteQueryWrite(ctx context.Context, query string, params map[string]any) ([]*neo4j.Record, error) {
	session := neo.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: neo.config.database})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, getTransactionFunc(ctx, query, params), neo4j.WithTxTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("could not execute neo4j session write: %w", err)
	}
	return result.([]*neo4j.Record), nil
}

func (neo Neo4jConn) PurgeDB(ctx context.Context) error {
	_, err := neo.ExecuteQueryWrite(ctx, queryPurgeDB, nil)
	if err != nil {
		return fmt.Errorf("could not purge neo4j db: %w", err)
	}

	// _, err = neo.ExecuteQueryWrite(ctx, createKeyConstraint, nil)
	// if err != nil {
	// 	return fmt.Errorf("could not create key constraint: %w", err)
	// }

	return nil
}

func getTransactionFunc(ctx context.Context, query string, params map[string]any) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		r, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, fmt.Errorf("could not execute neo4j transaction run: %w", err)
		}
		data, err := r.Collect(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not execute neo4j transaction collect: %w", err)
		}
		return data, nil
	}
}
