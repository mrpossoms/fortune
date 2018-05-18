package main

import (
	"math/rand"
	"math"
	// "fmt"
)

const (
	PlotSea float32 = 0
	PlotBeach       = 0.125
	PlotForest      = 0.5
	PlotMountain    = 1.0
)

type Plot struct {
	Explored int64
	Elevation float32
	Unit PlotUnit
	Tile PlotTile
}

type Box struct {
	X, Y int
	View struct {
		Width, Height int
	}
}

type Camera struct {
	Box
}

type World struct {
	Plots [200][200]Plot
	Smoothness int

	Cursor Box
	Cam Camera
}


func (c *Box) Move(dx, dy int) {
	nx, ny := c.X + dx, c.Y + dy

	if nx >= 0 && nx < 200 {
		c.X = nx
	}

	if ny >= 0 && ny < 200 {
		c.Y = ny
	}
}


func (c *Camera) Move(dx, dy int) {
	hw, hh := c.View.Width >> 1, c.View.Height >> 1
	nx, ny := c.X + dx, c.Y + dy

	if nx - hw >= 0 && nx + hw < 200 {
		c.X = nx
	}

	if ny - hh >= 0 && ny + hh < 200 {
		c.Y = ny
	}
}


func (p *Plot) productionRate() float32 {
	switch {
	case p.Elevation <= PlotSea:
		return 1
	case p.Elevation <= PlotBeach:
		return 0.1
	case p.Elevation <= PlotForest:
		return 2
	}

	return 0;
}


func (w *World) avgPatch(cx, cy int) float32 {
	avg := float32(0.0)
	n := 0
	for x := cx - 1; x <= cx + 1; x += 1 {
		for y := cy - 1; y <= cy + 1; y += 1 {
			avg += w.Plots[x][y].Elevation
			n += 1
		}
	}

	return avg / float32(n)
}

func (w *World) Init(seed int64) {
	r := rand.New(rand.NewSource(seed))

	w.Cursor.X = 100
	w.Cursor.Y = 100
	w.Cam.X = 100
	w.Cam.Y = 100
	w.Cam.View.Width, w.Cam.View.Height = 80, 32

	// populate with initial noize
	width, height := len(w.Plots), len(w.Plots[0])
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			plot := &w.Plots[x][y]

			dist := math.Sqrt(math.Pow(float64(x-w.Cam.X) / 2, 2) + math.Pow(float64(y-w.Cam.Y), 2))
			if dist <= 5 {
				plot.Explored = 1
			}

			if x == 0 || x == width - 1 || y == 0 || y == height - 1 {
				plot.Elevation = -1;
				continue
			}

			w := Gauss2D(float32(x) / float32(width),
			            float32(y) / float32(height),
						r.Float32(),
						r.Float32(), 0.5, 0.4) - 0.25
			plot.Elevation += float32(r.NormFloat64()) + w
		}
	}

	for i := w.Smoothness; i > 0; i -= 1 {
		for x := 1; x < len(w.Plots) - 1; x += 1 {
			for y := 1; y < len(w.Plots[x]) - 1; y += 1 {
				plot := &w.Plots[x][y]
				plot.Elevation = w.avgPatch(x, y)
				// fmt.Println(patch)
			}
		}
	}
}
