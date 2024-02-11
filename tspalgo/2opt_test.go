package tspalgo

import (
	"log"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	mylogger "tspgo/logger"
	"tspgo/tspinst"
)

var tspInst *tspinst.TspInst

func TestMain(m *testing.M) {
	mylogger.SetLevel(slog.LevelDebug)

	foldaPath := "../problems"
	fname := "a280.tsp"
	fpath := filepath.Join(foldaPath, fname)

	var err error
	tspInst, err = tspinst.LoadTspInst(fpath)
	if err != nil {
		mylogger.L().Error("Failed to load TSP instance", "err", err)
		panic(err)
	}

	mylogger.L().Info("Loading TSP instance is succeeded", "tspName", tspInst.ProblemName, "tspDim", tspInst.PointsDim)

	ret := m.Run()
	os.Exit(ret)
}

func getRand() *rand.Rand {
	seed := 1000
	resouce := rand.NewSource(int64(seed))
	return rand.New(resouce)
}

func Test2Opt(t *testing.T) {
	ans := ConstructWithGreedingByAllOriginConcurrently(nil, tspInst)

	randInst := getRand()
	n := ans.Inst.PointsDim

	for i := 0; i < 1000; i++ {
		idx1 := AnswerIndexT(randInst.Intn(n))
		idx2 := AnswerIndexT(randInst.Intn(n))

		Do2Opt(ans, idx1, idx2)
		if ok, msg := ans.IsCorrectAnswer(); !ok {
			mylogger.L().Error("Invalid answer", "idx1", idx1, "idx2", idx2, "msg", msg, "ans", ans.String())
			log.Fatalln("Invalid answer")
		} else {
			mylogger.L().Info("Valid answer", "idx1", idx1, "idx2", idx2, "i", i)
		}
	}
}

func TestCalc2OptScore(t *testing.T) {
	ans := ConstructWithGreedingByAllOriginConcurrently(nil, tspInst)
	ans.CalcScore()

	randInst := getRand()
	n := ans.Inst.PointsDim

	for i := 0; i < 1000; i++ {
		idx1 := AnswerIndexT(randInst.Intn(n))
		idx2 := AnswerIndexT(randInst.Intn(n))

		Do2Opt(ans, idx1, idx2)
		score2opt := ans.Score
		ans.CalcScore()

		if score2opt != ans.Score {
			mylogger.L().Error("Invalid score", "idx1", idx1, "idx2", idx2, "score2opt", score2opt, "score", ans.Score)
			log.Fatalln("Invalid score")
		} else {
			mylogger.L().Info("Valid score", "idx1", idx1, "idx2", idx2, "score", score2opt)

		}
	}
}
