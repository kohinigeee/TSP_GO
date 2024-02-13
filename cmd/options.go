package cmd

import (
	"flag"
	"fmt"
	"path/filepath"
)

type MainOptions struct {
	TspName   string
	TraceName string
}

func newMainOptions() *MainOptions {
	return &MainOptions{
		TspName:   "",
		TraceName: "",
	}
}

func LoadOptions() (*MainOptions, error) {
	options := &MainOptions{}

	flag.StringVar(&options.TspName, "tspName", "", "TSP file name")
	flag.StringVar(&options.TraceName, "traceName", "", "Trace file name")

	flag.Parse()

	if options.TspName == "" {
		flag.Usage()
		return nil, fmt.Errorf("TSP file name is not specified")
	}

	if options.TraceName == "" {
		flag.Usage()
		return nil, fmt.Errorf("Trace file name is not specified")
	}

	return options, nil
}

func (o *MainOptions) TspPath() string {
	problemFolda := "./problems"
	return filepath.Join(problemFolda, o.TspName+".tsp")
}

func (o *MainOptions) TracePath() string {
	traceFolda := "./trace"
	return filepath.Join(traceFolda, o.TraceName+".out")
}
