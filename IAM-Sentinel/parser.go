package main

import (
	"encoding/json"
)

type Parser struct{}

func NewParser() (*Parser, error) {
	p := Parser{}
	return &p, nil
}

func (p *Parser) parse(data []byte) (*GAAD, error) {
	var gaad GAAD

	if err := json.Unmarshal(data, &gaad); err != nil {
		return nil, err
	}

	return &gaad, nil
}
