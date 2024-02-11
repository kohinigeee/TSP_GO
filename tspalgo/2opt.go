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
