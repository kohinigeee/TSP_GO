package cmd

import (
	"context"
	"tspgo/tspalgo"
	"tspgo/tspinst"
)

func MainCtest() {
	options, err := LoadOptions()
	if err != nil {
		logger.Error("Failed to load options", "err", err)
		return
	}

	tspInst, err := tspinst.LoadTspInst(options.TspPath())
	if err != nil {
		logger.Error("Failed to load TSP instance", "err", err)
		return
	}

	logger.Info("Loading TSP instance is succeeded", "tspName", tspInst.ProblemName, "tspDim", tspInst.PointsDim)

	//初期解の構築

	ans := tspalgo.ConstructWithGreedingByAllOriginConcurrently(context.Background(), tspInst)
	ans.CalcScore()

	logger.Info("Initial solution is constructed", "ans", ans.String(), "score", ans.Score)
}
