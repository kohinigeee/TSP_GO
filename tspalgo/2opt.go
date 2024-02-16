package tspalgo

import (
	mylogger "tspgo/logger"
	"tspgo/tspinst"
	"tspgo/utility/slicecol"
)

// 2Opt後のスコアの増加量差分を利用して計算する
// next1-idex2, idex1-next2をnext1-next2, idex1-idex2に結び変える
func Calc2OptScoreDif(ans *TspAnswer, idx1, idx2 AnswerIndexT) tspinst.EcDistance {
	if !ans.isValidOrder(idx1) || !ans.isValidOrder(idx2) {
		mylogger.L().Warn("[CalcFast2OptScore] Invalid order", "idx1", idx1, "idx2", idx2)
		return 0
	}

	if idx1 == idx2 {
		mylogger.AL().Warn("[CalcFast2OptScore] idx1 and idx2 are same idx", "idx1", idx1, "idx2", idx2)
		return 0
	}

	next1 := ans.NextOrder(idx1)
	next2 := ans.NextOrder(idx2)

	if next1 == idx2 || next2 == idx1 {
		mylogger.AL().Warn("[CalcFast2OptScore] idx1 and idx2 are next", "idx1", idx1, "idx2", idx2)
		return 0
	}

	point1 := ans.Point(idx1)
	point2 := ans.Point(idx2)

	nextPoint1 := ans.Point(next1)
	nextPoint2 := ans.Point(next2)

	// 2Opt前のスコア
	beforeScore := point1.Distance(nextPoint1) + point2.Distance(nextPoint2)

	// 2Opt後のスコア
	afterScore := point1.Distance(point2) + nextPoint1.Distance(nextPoint2)

	return afterScore - beforeScore
}

// 2Opt操作を実行する
// next1-idex2, idex1-next2をnext1-next2, idex1-idex2に結び変える
func Do2Opt(ans *TspAnswer, idx1, idx2 AnswerIndexT) {
	if !ans.isValidOrder(idx1) || !ans.isValidOrder(idx2) {
		mylogger.L().Warn("[TspAnswer.Do2Opt] Invalid order", "idx1", idx1, "idx2", idx2)
		return
	}

	if idx1 == idx2 {
		mylogger.AL().Warn("[TspAnswer.Do2Opt] idx1 and idx2 are same idx", "idx1", idx1, "idx2", idx2)
		return
	}

	next1 := ans.NextOrder(idx1)
	next2 := ans.NextOrder(idx2)

	if next1 == idx2 || next2 == idx1 {
		mylogger.AL().Warn("[TspAnswer.Do2Opt] idx1 and idx2 are next, so cant't 2Opt", "idx1", idx1, "idx2", idx2)
		return
	}

	if !ans.IsCalculatedScore() {
		ans.CalcScore()
	}

	scoreDif := Calc2OptScoreDif(ans, idx1, idx2)
	ans.Score += scoreDif

	if next1 < idx2 {
		ans.Order = slicecol.Reverse[tspinst.PointIndexT](ans.Order, int(next1), int(idx2))
	} else {
		ans.Order = slicecol.Reverse[tspinst.PointIndexT](ans.Order, int(next2), int(idx1))
	}
}

type AnswerMoveFunc func(ans *TspAnswer) bool

func MoveBestNeighborBy2Opt(ans *TspAnswer) (isMoved bool) {
	if !ans.IsCalculatedScore() {
		ans.CalcScore()
	}

	bestDif := tspinst.EcDistance(0)
	bestIdx1 := AnswerIndexT(-1)
	bestIdx2 := AnswerIndexT(-1)

	for i := 0; i < ans.Inst.PointsDim; i++ {
		for j := i + 1; j < ans.Inst.PointsDim; j++ {
			idx1 := AnswerIndexT(i)
			idx2 := AnswerIndexT(j)
			dif := Calc2OptScoreDif(ans, idx1, idx2)

			if dif < bestDif {
				bestDif = dif
				bestIdx1 = idx1
				bestIdx2 = idx2
			}
		}
	}

	if bestDif == 0 {
		isMoved = false
		return
	}

	Do2Opt(ans, bestIdx1, bestIdx2)
	isMoved = true
	return
}

func MoveBestNeighborBy2OptConcurrently(ans *TspAnswer) (isMoved bool) {
	if !ans.IsCalculatedScore() {
		ans.CalcScore()
	}

	type DifInfo struct {
		idx1 AnswerIndexT
		idx2 AnswerIndexT
		dif  tspinst.EcDistance
	}

	const initialDif = tspinst.EcDistance(0)
	bestDifInfo := DifInfo{idx1: -1, idx2: -1, dif: initialDif}

	difInfoCh := make(chan DifInfo)

	for i := 0; i < ans.Inst.PointsDim; i++ {
		go func(targetIdx int) {
			tmpBestDifInfo := DifInfo{idx1: -1, idx2: -1, dif: initialDif}
			idx1 := AnswerIndexT(targetIdx)
			for j := 0; j < ans.Inst.PointsDim; j++ {
				idx2 := AnswerIndexT(j)
				dif := Calc2OptScoreDif(ans, idx1, idx2)
				if dif < tmpBestDifInfo.dif {
					tmpBestDifInfo = DifInfo{idx1: idx1, idx2: idx2, dif: dif}
				}
			}

			difInfoCh <- tmpBestDifInfo
		}(i)
	}

	compDifInfoIdx := func(info1, info2 *DifInfo) bool {
		if info1.idx1 != info2.idx1 {
			return info1.idx1 < info2.idx1
		}
		return info1.idx2 < info2.idx2
	}

	for i := 0; i < ans.Inst.PointsDim; i++ {
		tmpDifInfo := <-difInfoCh
		if tmpDifInfo.dif >= initialDif {
			continue
		}

		// 差分が同じ場合はidx1, idx2の値が小さい方を優先する
		//(直列実行と同じ結果になるようにするため)
		if tmpDifInfo.dif == bestDifInfo.dif {
			if compDifInfoIdx(&tmpDifInfo, &bestDifInfo) {
				bestDifInfo = tmpDifInfo
			}
			continue
		}

		if tmpDifInfo.dif < bestDifInfo.dif {
			bestDifInfo = tmpDifInfo
		}
	}

	if bestDifInfo.dif == initialDif {
		isMoved = false
		return
	}

	Do2Opt(ans, bestDifInfo.idx1, bestDifInfo.idx2)
	isMoved = true
	return
}

func LocalSearchBy2Opt(ans *TspAnswer, moveFunc AnswerMoveFunc) (moveCount int) {
	moveCount = 0
	for {
		if !moveFunc(ans) {
			break
		}
		moveCount++
		mylogger.L().Debug("LocalSearchBy2Opt", "moveCount", moveCount, "score", ans.Score)
	}

	return
}
