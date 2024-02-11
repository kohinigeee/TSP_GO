package tspalgo

import (
	"context"
	"fmt"
	"runtime/trace"
	"tspgo/tspinst"
	"tspgo/utility/slicecol"
)

func ConstructWithGreedingByOneOrigin(ctx context.Context, inst *tspinst.TspInst, origin tspinst.PointIndexT) *TspAnswer {

	if ctx != nil {
		region := trace.StartRegion(ctx, fmt.Sprintf("Region: ConstructWithGreeding (PloblemName: %v, Origin: %v)", inst.ProblemName, origin))
		region.End()
	}

	unSelecetedPointIndexs := make([]tspinst.PointIndexT, inst.PointsDim)
	for i := 0; i < inst.PointsDim; i++ {
		unSelecetedPointIndexs[i] = tspinst.PointIndexT(i)
	}

	unSelecetedPointIndexs = slicecol.RemoveFast[tspinst.PointIndexT](unSelecetedPointIndexs, int(origin))

	ansOrder := make([]tspinst.PointIndexT, inst.PointsDim)
	ansOrder[0] = origin

	for i := 1; i < inst.PointsDim; i++ {
		tailOrder := ansOrder[i-1]
		idx, nearestPointOrder := CalcNearestPoint(inst, tailOrder, unSelecetedPointIndexs)
		ansOrder[i] = nearestPointOrder
		unSelecetedPointIndexs = slicecol.RemoveFast[tspinst.PointIndexT](unSelecetedPointIndexs, idx)
	}

	return NewTspAnswer(inst, ansOrder)
}

func ConstructWithGreedingByAllOriginConcurrently(ctx context.Context, inst *tspinst.TspInst) *TspAnswer {
	var bestAns *TspAnswer = nil

	if ctx != nil {
		region := trace.StartRegion(ctx, fmt.Sprintf("Region: ConstructWithGreedingByAllOriginConcurrently (PloblemName: %v)", inst.ProblemName))
		defer region.End()
	}

	ch := make(chan *TspAnswer)
	defer close(ch)

	for i := 0; i < inst.PointsDim; i++ {
		go func(i int) {
			ans := ConstructWithGreedingByOneOrigin(nil, inst, tspinst.PointIndexT(i))
			ans.CalcScore()
			ch <- ans
		}(i)
	}

	for i := 0; i < inst.PointsDim; i++ {
		tmpAns := <-ch
		if bestAns == nil {
			bestAns = tmpAns
		} else {
			if bestAns.Score > tmpAns.Score {
				bestAns = tmpAns
			}
		}
	}
	return bestAns
}

func ConstructWithGreedingByAllOrigin(ctx context.Context, inst *tspinst.TspInst) *TspAnswer {
	var bestAns *TspAnswer = nil

	if ctx != nil {
		region := trace.StartRegion(ctx, fmt.Sprintf("Region: ConstructWithGreedingByAllOrigin (PloblemName: %v)", inst.ProblemName))
		defer region.End()
	}

	bestAns = ConstructWithGreedingByOneOrigin(nil, inst, 0)
	bestAns.CalcScore()

	for i := 1; i < inst.PointsDim; i++ {
		tmpAns := ConstructWithGreedingByOneOrigin(nil, inst, tspinst.PointIndexT(i))
		tmpAns.CalcScore()

		if bestAns.Score > tmpAns.Score {
			bestAns = tmpAns
		}
	}

	return bestAns
}

// originから最も近い点を探す
// { 戻り値: 最も近い点の番号の添え字, 最も近い点のインデックス }
func CalcNearestPoint(inst *tspinst.TspInst, origin tspinst.PointIndexT, orderIdxs []tspinst.PointIndexT) (int, tspinst.PointIndexT) {
	if len(orderIdxs) == 0 {
		panic("[CalcNearestPoint] Invalid orderIdxs length")
	}

	if len(orderIdxs) == 1 {
		return 0, orderIdxs[0]
	}

	var nearestOrderIdx int = 0
	var nearestDist tspinst.EcDistance = inst.Point(origin).Distance(inst.Point(orderIdxs[nearestOrderIdx]))

	for i, order := range orderIdxs {
		if i == 0 {
			continue
		}

		dist := inst.Point(origin).Distance(inst.Point(order))
		if dist < nearestDist {
			nearestDist = dist
			nearestOrderIdx = i
		}
	}

	return nearestOrderIdx, orderIdxs[nearestOrderIdx]
}
