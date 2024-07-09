package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type Exporter struct {
	logger     *logrus.Logger
	exportPath string
}

func NewExporter(logger *logrus.Logger, exportPath string) (*Exporter, error) {
	e := Exporter{logger: logger, exportPath: exportPath}

	dir := path.Dir(exportPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	return &e, nil
}

func (e *Exporter) logNodes(nodes []CentralityNode) {
	for _, node := range nodes {
		jsonBytes, err := json.Marshal(node)
		if err != nil {
			e.logger.Error(err)
			continue
		}
		e.logger.Info(string(jsonBytes))
	}
}

func (e *Exporter) writeToFile(nodes []CentralityNode, fileName string) {
	jsonBytes, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		e.logger.Error(err)
		return
	}

	f, err := os.Create(path.Join(e.exportPath, fileName))
	if err != nil {
		e.logger.Error(err)
		return
	}
	defer f.Close()

	n, err := f.Write(jsonBytes)
	if err != nil {
		e.logger.Error(err)
	}

	if n == 0 {
		e.logger.Warn("no bytes were written to file", "filename", fileName)
	}
}
