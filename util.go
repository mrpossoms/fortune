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

func Gauss2D(x, y, muX, muY, sigX, sigY float32) float32 {
	deno := float64(sigX + sigY) * math.Sqrt(2.0 * math.Pi)

	num_x := -math.Pow(float64((x - muX) / sigX), 2.0)
	num_y := -math.Pow(float64((y - muY) / sigY), 2.0)
	exp := (num_x + num_y) / 2.0

	return float32(math.Pow(math.E, exp) / deno)
}