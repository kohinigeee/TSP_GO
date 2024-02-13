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

	wholeTimer := exetimer.MeasureStart()

	//-----------直列実行-----------
	ctx, normalTask := trace.NewTask(context.Background(), "MainSearchSequentaly")

	//初期解の構築
	initialTimer := exetimer.MeasureStart()
	region1 := trace.StartRegion(ctx, "Region: ConstructWithGreedingByAllOrigin")
	ans := tspalgo.ConstructWithGreedingByAllOrigin(context.Background(), tspInst)
	ans.CalcScore()
	region1.End()
	initialTimer.MeasureEnd()

	cpAns1 := ans.Copy()
	logger.Info("Constructing with greeding was finished", "time(ms)", initialTimer.ElapsedMilliSeconds(), "ans", ans.String())

	//局所探索
	localSearchTimer := exetimer.MeasureStart()
	region2 := trace.StartRegion(ctx, "Region: LocalSearch")
	seqMoveCnt := tspalgo.LocalSearchBy2Opt(ans, tspalgo.MoveBestNeighborBy2Opt)
	region2.End()
	localSearchTimer.MeasureEnd()

	cpAns2 := ans.Copy()
	logger.Info("Local search was finished", "time(ms)", localSearchTimer.ElapsedMilliSeconds(), "ans", ans.String())

	normalTask.End()

	seqWholeTime := localSearchTimer.ElapsedMilliSeconds() + initialTimer.ElapsedMilliSeconds()

	//-----------並列実行-----------
	ctx, concurrentTask := trace.NewTask(context.Background(), "MainSearchConcurrently")

	conInitialTimer := exetimer.MeasureStart()
	region1 = trace.StartRegion(ctx, "Region: ConstructWithGreedingByAllOriginConcurrently")
	ans = tspalgo.ConstructWithGreedingByAllOriginConcurrently(context.Background(), tspInst)
	ans.CalcScore()
	region1.End()
	conInitialTimer.MeasureEnd()

	cpCoInitislAns := ans.Copy()
	logger.Info("Constructing with concurrently greeding was finished", "time(ms)", conInitialTimer.ElapsedMilliSeconds(), "ans", ans.String())

	//局所探索
	conLocalSearchTimer := exetimer.MeasureStart()
	region2 = trace.StartRegion(ctx, "Region: LocalSearchConcurrently")
	conMoveCnt := tspalgo.LocalSearchBy2Opt(ans, tspalgo.MoveBestNeighborBy2OptConcurrently)
	region2.End()
	conLocalSearchTimer.MeasureEnd()

	cpCoLsearchAns := ans.Copy()
	logger.Info("Concurrently local search was finished", "time(ms)", conLocalSearchTimer.ElapsedMilliSeconds(), "ans", ans.String())

	concurrentTask.End()
	conWholeTime := conInitialTimer.ElapsedMilliSeconds() + conLocalSearchTimer.ElapsedMilliSeconds()

	wholeTimer.MeasureEnd()

	//-----------結果の出力-----------
	logger.Info("Whole process was finished", "time(ms)", wholeTimer.ElapsedMilliSeconds())

	logger.Info("--------------Sequentaly-----------------------")

	logger.Info("[1] Sequentaly constructing with greeding was finished", "time(ms)", initialTimer.ElapsedMilliSeconds(), "score", cpAns1.Score, "ans", cpAns1.String())
	logger.Info("[2] Sequentaly Local search was finished", "time(ms)", localSearchTimer.ElapsedMilliSeconds(), "score", cpAns2.Score, "ans", cpAns2.String(), "moveCnt", seqMoveCnt)
	logger.Info("[Whole] Sequentaly whole time", "time(ms)", seqWholeTime)

	logger.Info("--------------Concurrently-----------------------")
	logger.Info("[1] Concurrently constructing with greeding was finished", "time(ms)", conInitialTimer.ElapsedMilliSeconds(), "score", cpCoInitislAns.Score, "ans", cpCoInitislAns.String())
	logger.Info("[2] Concurrently Local search was finished", "time(ms)", conLocalSearchTimer.ElapsedMilliSeconds(), "score", cpCoLsearchAns.Score, "ans", cpCoLsearchAns.String(), "moveCnt", conMoveCnt)
	logger.Info("[Whole] Concurrently whole time", "time(ms)", conWholeTime)
}
