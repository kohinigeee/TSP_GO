package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
	mylogger "tspgo/logger"
	"tspgo/tspalgo"
	"tspgo/tspinst"
	"tspgo/utility/exetimer"
)

func MainWithGreeding() {
	logger := mylogger.L()

	var fname string
	var tracefname string
	foldaPath := "./problems"
	flag.StringVar(&fname, "tspName", "", "TSP file name")
	flag.StringVar(&tracefname, "traceName", "", "Trace file name")

	flag.Parse()

	if fname == "" {
		logger.Error("TSP file name is not specified")
		fmt.Println()
		flag.Usage()
		return
	}
	fname = fname + ".tsp"
	fpath := filepath.Join(foldaPath, fname)

	tspInst, err := tspinst.LoadTspInst(fpath)
	if err != nil {
		logger.Error("Failed to load TSP instance", "err", err)
		return
	}

	logger.Info("Loading TSP instance is succeeded", "tspName", tspInst.ProblemName, "tspDim", tspInst.PointsDim)

	if tracefname == "" {
		logger.Error("Trace file name is not specified")
		fmt.Println()
		flag.Usage()
	}

	tracefname = tracefname + ".out"
	traceFoldaPath := "./trace"
	traceFilepath := filepath.Join(traceFoldaPath, tracefname)

	traceFile, err := os.Create(traceFilepath)
	if err != nil {
		log.Fatalln("Failed to create trace file", "err", err)
	}

	defer func() {
		if err := traceFile.Close(); err != nil {
			log.Fatalln("Failed to close trace file", "err", err)
		}
	}()

	if err := trace.Start(traceFile); err != nil {
		log.Fatalln("Failed to start trace", "err", err)
	}
	defer trace.Stop()

	ctx, task := trace.NewTask(context.Background(), "MainWithGreeding1")
	timer1 := exetimer.MeasureStart()
	ans1 := tspalgo.ConstructWithGreedingByAllOrigin(ctx, tspInst)
	logger.Info("Calculated Sequentially")
	ans1.CalcScore()
	timer1.MeasureEnd()
	trace.Log(ctx, "Calculated Sequentially", "Message1")
	task.End()

	ctx, task = trace.NewTask(context.Background(), "MainWithGreeding2")
	timer2 := exetimer.MeasureStart()
	ans2 := tspalgo.ConstructWithGreedingByAllOriginConcurrently(ctx, tspInst)
	logger.Info("Calculated Concurrently")
	ans2.CalcScore()
	timer2.MeasureEnd()
	task.End()

	logger.Info("Answer1", "time(ms)", timer1.ElapsedMilliSeconds(), "ans1", ans1)
	logger.Info("Answer2", "time(ms)", timer2.ElapsedMilliSeconds(), "ans2", ans2)
}

func MainWithPprof() {
	logger := mylogger.L()

	var fname string
	foldaPath := "./problems"
	flag.StringVar(&fname, "tspName", "", "TSP file name")

	flag.Parse()

	if fname == "" {
		logger.Error("TSP file name is not specified")
		fmt.Println()
		flag.Usage()
		return
	}
	fname = fname + ".tsp"
	fpath := filepath.Join(foldaPath, fname)

	tspInst, err := tspinst.LoadTspInst(fpath)
	if err != nil {
		logger.Error("Failed to load TSP instance", "err", err)
		return
	}

	logger.Info("Loading TSP instance is succeeded", "tspName", tspInst.ProblemName, "tspDim", tspInst.PointsDim)

	ppf, err := os.Create("cpu2.pprof")
	if err != nil {
		logger.Error("Failed to create pprof file", "err", err)
		return
	}
	defer ppf.Close()

	if err := pprof.StartCPUProfile(ppf); err != nil {
		logger.Error("Failed to start CPU profile", "err", err)
		return
	}

	ans := tspalgo.ConstructWithGreedingByAllOriginConcurrently(nil, tspInst)
	ans.CalcScore()
	logger.Info("Answer", "ans", ans)

	pprof.StopCPUProfile()
}
