package tspinst

import "math"

type EcDistance int

type Point struct {
	X, Y          int
	OriginalIndex int
}

func NewPoint(x, y, idx int) *Point {
	return &Point{
		X:             x,
		Y:             y,
		OriginalIndex: idx,
	}
}

// 2点間のユークリッド距離を求める
func (p *Point) Distance(q *Point) EcDistance {
	dx := p.X - q.X
	dy := p.Y - q.Y
	dec := math.Sqrt(float64(dx*dx + dy*dy))
	return EcDistance(dec + 0.5)
}
