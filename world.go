package main

import (
	"math/rand"
	// "math"
	"fmt"
	// "net"
)

const (
	PlotSea float32 = 0
	PlotBeach       = 0.125
	PlotPlains      = 0.5
	PlotForest      = 0.8
	PlotMountain    = 4.0
)

const (
	WorldWidth = 200
	WorldHeight = 200
	ViewWidth = 80
	ViewHeight = 32
)

type Plot struct {
	Explored int64
	Elevation float32
	Unit PlotUnit
	Tile PlotTile
	X, Y int
}

type World struct {
	Plots [WorldWidth][WorldHeight]Plot
	Smoothness int

	R *rand.Rand
}


func (p *Plot) TerrainName() string {
	switch {
	case p.Elevation < PlotSea:
		return "sea"
	case  p.Elevation < PlotBeach:
		return "beach"
	case p.Elevation < PlotPlains:
		return "plains"
	case  p.Elevation < PlotForest:
		return "forest"
	case p.Elevation < PlotMountain:
		return "mountains"
	}

	return ""
}


func (p *Plot) Description(descType int) string {
	desc := fmt.Sprintf("%v (%d,%d)", p.TerrainName(), p.X, p.Y)

	if p.Unit.Type != UnitNone {
		return p.Unit.Description(descType) + " in the " + desc
	}

	return desc
}


func (p *Plot) SpawnUnit(unitType int, owner *Player) (*PlotUnit, string) {
	if p.Unit.Type != UnitNone {
		return nil, "Couldn't spawn " + Units[unitType].Name + " in occupied plot!"
	}
	p.Unit = Units[unitType]
	p.Unit.OwnerID = owner.ID

	return &p.Unit, fmt.Sprintf("%v spawned!", Units[unitType].Name)
}


func (p *Plot) ProductionRate() float32 {
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


func (p *Plot) IsLivable() bool {
	return p.Elevation > PlotSea && p.Elevation < PlotForest
}


func (w *World) FindLivablePlot() *Plot {
	for {
		x, y := w.R.Intn(WorldWidth), w.R.Intn(WorldHeight)
		plot := &w.Plots[x][y]
		if plot.IsLivable() {
			return plot
		}
	}
}


func (w *World) avgPatchElevation(cx, cy int) float32 {
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

func (w *World) Reveal(x, y, r int, playerId int64) (min_x, min_y, max_x, max_y int) {
	min_x, max_x = max(x - r, 0), min(x + r, len(w.Plots)-1)
	min_y, max_y = max(y - r / 2, 0), min(y + r / 2, len(w.Plots[0])-1)

	for i := min_x; i <= max_x; i += 1 {
		for j := min_y; j <= max_y; j += 1 {
			// dx, dy := x - i, y - j
			// dist := (dx * dx) + (dy * dy)//math.Sqrt(math.Pow(float64(x-i) / 2, 2) + math.Pow(float64(y-j), 2))
			// if dist <= (r * r) {
			// 	w.Plots[i][j].Explored |= playerId
			// }

			// w.Plots[i][j].Explored = 0

			w.Plots[i][j].Explored |= playerId
		}
	}

	// min_x, min_y = 0, 0
	// max_x, max_y = 199, 199

	return
}

func (w *World) Init(seed int64) {
	if w.R == nil {
		w.R = rand.New(rand.NewSource(seed))
	}

	// populate with initial noise
	width, height := len(w.Plots), len(w.Plots[0])
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			plot := &w.Plots[x][y]

			plot.X, plot.Y = x, y
			plot.Explored = 0

			if x == 0 || x == width - 1 || y == 0 || y == height - 1 {
				plot.Elevation = -1;
				continue
			}

			off := Gauss2D(float32(x) / float32(width),
			               float32(y) / float32(height),
			               w.R.Float32(),
			               w.R.Float32(), 0.5, 0.4) - 0.25
			plot.Elevation += float32(w.R.NormFloat64()) + off
		}
	}

	for i := w.Smoothness; i > 0; i -= 1 {
		for x := 1; x < len(w.Plots) - 1; x += 1 {
			for y := 1; y < len(w.Plots[x]) - 1; y += 1 {
				plot := &w.Plots[x][y]
				plot.Elevation = w.avgPatchElevation(x, y)
				// fmt.Println(patch)
			}
		}
	}
}
