package main

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Parser struct {
	logger *logrus.Logger
}

func NewParser(logger *logrus.Logger) (*Parser, error) {
	p := Parser{logger: logger}
	return &p, nil
}

func (p *Parser) parse(data []byte) (*GAAD, error) {
	var gaad GAAD

	if err := json.Unmarshal(data, &gaad); err != nil {
		return nil, fmt.Errorf("could not unmarshall data into GAAD: %w", err)
	}

	return &gaad, nil
}
