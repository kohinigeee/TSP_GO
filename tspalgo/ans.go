package tspalgo

import (
	"fmt"
	"math"
	mylogger "tspgo/logger"
	"tspgo/tspinst"
)

const (
	NULL_SCORE                    = -1
	INF_SCORE  tspinst.EcDistance = math.MaxInt64
)

// 解の点の順番を表す型
type AnswerIndexT int

type TspAnswer struct {
	Inst  *tspinst.TspInst
	Order []tspinst.PointIndexT
	Score tspinst.EcDistance
}

func NewTspAnswer(inst *tspinst.TspInst, order []tspinst.PointIndexT) *TspAnswer {
	if inst.PointsDim < 1 {
		panic("[NewTspAnswer] Invalid points dimension")
	}

	if len(order) != inst.PointsDim {
		panic("[NewTspAnswer] Invalid order length" + fmt.Sprintf(" len:%v, dim:%v", len(order), inst.PointsDim))
	}

	return &TspAnswer{
		Inst:  inst,
		Order: order,
		Score: NULL_SCORE,
	}
}

func (ans *TspAnswer) String() string {
	ansStr := ""
	ansStr += fmt.Sprintf("InstName: %v ", ans.Inst.ProblemName)
	if ans.IsCalculatedScore() {
		ansStr += fmt.Sprintf(", Score: %v", ans.Score)
	}
	ansStr += "\n"
	ansStr += "{ "
	for i, order := range ans.Order {
		if i > 70 {
			ansStr += fmt.Sprintf("...(%v)", len(ans.Order)-i)
			break
		}
		ansStr += fmt.Sprintf("%v ", order)
	}
	ansStr += "}"

	return ansStr
}

func (ans *TspAnswer) PointDim() int {
	return ans.Inst.PointsDim
}

func (ans *TspAnswer) TailPoint() *tspinst.Point {
	tailIndex := ans.Order[ans.PointDim()-1]
	return ans.Inst.Point(tailIndex)
}

func (ans *TspAnswer) Point(i AnswerIndexT) *tspinst.Point {
	return ans.Inst.Point(ans.Order[i])
}

func (ans *TspAnswer) NextOrder(i AnswerIndexT) AnswerIndexT {
	return (i + 1) % AnswerIndexT(ans.Inst.PointsDim)
}

func (ans *TspAnswer) PreviousOrder(i AnswerIndexT) AnswerIndexT {
	pointsDim := AnswerIndexT(ans.Inst.PointsDim)
	return (i + pointsDim - 1) % pointsDim
}

func (ans *TspAnswer) IsCalculatedScore() bool {
	return ans.Score != NULL_SCORE
}

func (ans *TspAnswer) CalcScore() tspinst.EcDistance {
	if ok, msg := ans.IsCorrectAnswer(); !ok {
		mylogger.L().Error("[TspAnswer.CalcScore] ans is incorrect answer", "msg", msg)
		panic("[TspAnswer.CalcScore] ans is incorrect answer")
	}

	var sum tspinst.EcDistance = 0
	for i := 0; i < ans.Inst.PointsDim; i++ {
		idx := AnswerIndexT(i)
		p1 := ans.Point(idx)
		p2 := ans.Point(ans.NextOrder(idx))
		sum += p1.Distance(p2)
	}

	ans.Score = sum
	return sum
}

func (ans *TspAnswer) IsCorrectAnswer() (bool, string) {
	if len(ans.Order) != ans.Inst.PointsDim {

		return false, fmt.Sprintf("<Invalid order length> lengh: %v, dim: %v", len(ans.Order), ans.Inst.PointsDim)
	}

	sets := make(map[tspinst.PointIndexT]bool)
	for _, order := range ans.Order {

		if order < 0 || order >= tspinst.PointIndexT(ans.Inst.PointsDim) {
			return false, fmt.Sprintf("<Invalid order> order: %v, dim: %v", order, ans.Inst.PointsDim)
		}

		if _, ok := sets[order]; ok {
			return false, fmt.Sprintf("<Duplicated order> order: %v", order)
		}

		sets[order] = true
	}

	return true, ""
}

func (ans *TspAnswer) isValidOrder(idx AnswerIndexT) bool {
	return (idx >= 0) && (idx < AnswerIndexT(ans.Inst.PointsDim))
}

func (ans *TspAnswer) Do2Opt(idx1, idx2 AnswerIndexT) {
	if !ans.isValidOrder(idx1) || !ans.isValidOrder(idx2) {
		panic("[TspAnswer.Do2Opt] Invalid order")
	}

	if idx1 == idx2 {
		mylogger.L().Warn("[TspAnswer.Do2Opt] idx1 and idx2 are same idx", "idx1", idx1, "idx2", idx2)
		return
	}

}
