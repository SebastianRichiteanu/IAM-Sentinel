package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type (
	Sentinel struct {
		logger   *logrus.Logger
		dbConn   *Neo4jConn
		parser   *Parser
		mapper   *ResourceMapper
		analyzer *Analyzer
		exporter *Exporter
	}

	SentinelConfig struct {
		neo4jConfig  *Neo4jConfig
		exportPath   string
		defaultLimit int
	}
)

func NewSentinel(ctx context.Context, cfg SentinelConfig) (*Sentinel, error) {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.InfoLevel

	logger.Info("Initializing...")

	dbConn, err := InitializeDB(ctx, *cfg.neo4jConfig)
	if err != nil {
		return nil, err
	}

	parser, err := NewParser(logger)
	if err != nil {
		return nil, err
	}

	mapper, err := NewResourceMapper(logger, dbConn, parser)
	if err != nil {
		return nil, err
	}

	analyzer, err := NewAnalyzer(logger, dbConn, cfg.defaultLimit)
	if err != nil {
		return nil, err
	}

	exporter, err := NewExporter(logger, cfg.exportPath)
	if err != nil {
		return nil, err
	}

	s := &Sentinel{
		logger:   logger,
		dbConn:   dbConn,
		parser:   parser,
		mapper:   mapper,
		analyzer: analyzer,
		exporter: exporter,
	}

	return s, nil
}

func (s *Sentinel) stop(ctx context.Context) {
	s.dbConn.Stop(ctx)
}
