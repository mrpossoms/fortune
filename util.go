package main

import (
	"math"
)

// func RandNorm(mu, sig float32) float32 {
// 	deno := sig * float32(math.Sqrt(2 * math.Pi))
// 	max_y := 1.0 / deno
//
//
//
// 	return 0
// }

type Point struct {
	X, Y int
}

type Region struct {
	Min Point
	Max Point
}


func (r *Region) Area() int {
	return (1 + r.Max.X - r.Min.X) * (1 + r.Max.Y - r.Min.Y)
}


func Gauss2D(x, y, muX, muY, sigX, sigY float32) float32 {
	deno := float64(sigX + sigY) * math.Sqrt(2.0 * math.Pi)

	num_x := -math.Pow(float64((x - muX) / sigX), 2.0)
	num_y := -math.Pow(float64((y - muY) / sigY), 2.0)
	exp := (num_x + num_y) / 2.0

	return float32(math.Pow(math.E, exp) / deno)
}


func abs(i int) int {
	if i < 0 { return -i }
	return i
}

func max(a, b int) int {
	if a > b { return a }
	return b
}

func min(a, b int) int {
	if a < b { return a }
	return b
}
