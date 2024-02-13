package cmd

import (
	"context"
	"os"
	"runtime/trace"
	mylogger "tspgo/logger"
	"tspgo/tspalgo"
	"tspgo/tspinst"
	"tspgo/utility/exetimer"
)

var (
	logger = mylogger.L()
)

func MainWithLocalSearch() {

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

	traceFile, err := os.Create(options.TracePath())
	if err != nil {
		logger.Error("Failed to create trace file", "err", err)
		return
	}

	defer func() {
		if err := traceFile.Close(); err != nil {
			logger.Error("Failed to close trace file", "err", err)
			return
		}
	}()

	if err := trace.Start(traceFile); err != nil {
		logger.Error("Failed to start trace", "err", err)
		return
	}
	defer trace.Stop()

	ctx, task := trace.NewTask(context.Background(), "MainWithLocalSearch")

	wholeTimer := exetimer.MeasureStart()

	//初期解の構築
	initialTimer := exetimer.MeasureStart()
	region1 := trace.StartRegion(ctx, "Region: ConstructWithGreedingByAllOrigin")
	ans := tspalgo.ConstructWithGreedingByAllOriginConcurrently(context.Background(), tspInst)
	ans.CalcScore()
	region1.End()
	initialTimer.MeasureEnd()

	cpAns1 := ans.Copy()
	logger.Info("Constructing with greeding was finished", "time(ms)", initialTimer.ElapsedMilliSeconds(), "ans", ans.String())

	//局所探索
	localSearchTimer := exetimer.MeasureStart()
	region2 := trace.StartRegion(ctx, "Region: LocalSearch")
	tspalgo.LocalSearchBy2Opt(ans, tspalgo.MoveBestNeighborBy2OptConcurrently)
	// tspalgo.LocalSearchBy2Opt(ans, tspalgo.MoveBestNeighborBy2Opt)
	region2.End()
	localSearchTimer.MeasureEnd()

	wholeTimer.MeasureEnd()

	cpAns2 := ans.Copy()
	logger.Info("Local search was finished", "time(ms)", localSearchTimer.ElapsedMilliSeconds(), "ans", ans.String())

	task.End()

	logger.Info("Whole process was finished", "time(ms)", wholeTimer.ElapsedMilliSeconds(), "ans", ans.String())
	logger.Info("[1] Constructing with greeding was finished", "time(ms)", initialTimer.ElapsedMilliSeconds(), "score", cpAns1.Score, "ans", cpAns1.String())
	logger.Info("[2] Local search was finished", "time(ms)", localSearchTimer.ElapsedMilliSeconds(), "score", cpAns2.Score, "ans", cpAns2.String())
}
