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
	WorldWidth = 100
	WorldHeight = 50
	ViewWidth = 80
	ViewHeight = 32
)

type Plot struct {
	Explored int64
	Elevation float32
	Unit PlotUnit
	Tile PlotTile
	X, Y int
	Updated int
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


func (p *Plot) Neighbors(world *World) []*Plot {
	slice := make([]*Plot, 0, 8)

	for i := -1; i <= 1; i += 1 {
		for j := -1; j <= 1; j += 1 {
			r, c := p.X + i, p.Y + j
			if i | j == 0 { continue }
			if r < 0 || r > WorldWidth { continue }
			if c < 0 || c > WorldHeight { continue }
			n := &world.Plots[r][c]
			// GfxMsg(fmt.Sprintf("(%d, %d) -> (%d, %d) unit -> %d", r, c, n.X, n.Y, n.Unit.Type))

			slice = append(slice, n)
		}
	}

	return slice
}


func (p *Plot) HasNeighbor(world *World, neighbors []*Plot, unitType int, owner int64) bool {
	for ni := 0; ni < len(neighbors); ni += 1 {
		neighbor := neighbors[ni]
		unit := neighbor.Unit
		if unit.Type == unitType {
			return true
		}
	}

	return false
}


func (p *Plot) PossibleBuilds(world *World, owner int64) []int {
	unitCity := UnitIndex("city")
	unitVillage := UnitIndex("village")
	unitFarm := UnitIndex("farm")
	unitMine := UnitIndex("mine")

	neighbors := p.Neighbors(world)
	nextToVillage := p.HasNeighbor(world, neighbors, unitVillage, owner)
	nextToCity := p.HasNeighbor(world, neighbors, unitCity, owner)
	nextToFarm := p.HasNeighbor(world, neighbors, unitFarm, owner)

	buildables := make([]int, 0, 10)

	// TODO: figure out why this causes problems on the
	// client, but not the server
	if p.Unit.Type != UnitIndex("vacant") {
		return buildables
	}

	switch {
	case p.Elevation < PlotSea:
		break
	case  p.Elevation < PlotBeach:
		if nextToVillage || nextToCity {
			buildables = append(buildables, unitCity)
			buildables = append(buildables, unitVillage)
		}
		break
	case p.Elevation < PlotPlains:
		if nextToVillage || nextToCity || nextToFarm {
			buildables = append(buildables, unitCity)
			buildables = append(buildables, unitVillage)
			buildables = append(buildables, unitFarm)
		}
		break
	case  p.Elevation < PlotForest:
		if nextToVillage || nextToCity || nextToFarm {
			buildables = append(buildables, unitCity)
			buildables = append(buildables, unitVillage)
			buildables = append(buildables, unitFarm)
		}
	case p.Elevation < PlotMountain:
		if nextToVillage || nextToCity || nextToFarm {
			buildables = append(buildables, unitMine)
		}
	}

	return buildables
}


func (p *Plot) Description(descType int) string {
	desc := fmt.Sprintf("%v (%d,%d)", p.TerrainName(), p.X, p.Y)

	if p.Unit.Type != UnitIndex("vacant") {
		return p.Unit.Description(descType) + " in the " + desc
	}

	return desc
}


func (p *Plot) SpawnUnit(unitType int, owner *Player) (*PlotUnit, string) {
	if p.Unit.Type != UnitIndex("vacant") {
		return nil, "Couldn't spawn " + Units[unitType].Name + " in occupied plot!"
	}
	p.Unit = Units[unitType]
	p.Unit.OwnerID = owner.ID
	msg := fmt.Sprintf("(%d,%d) %v spawned for %v", p.X, p.Y, p.Unit.Name, owner.Name)

	fmt.Println(msg)

	return &p.Unit, msg
}


func (p *Plot) BuildMenu(world *World, owner int64, onSelection func(int)) {
	buildables := p.PossibleBuilds(world, owner)
	options := [10]string{}
	optSlice := options[0:0]

	for i := 0; i < len(buildables); i += 1 {
		ui := buildables[i]
		optSlice = append(optSlice, Units[ui].Name)
	}

	GfxMenu("These units can be built here", optSlice, func(selection int) {
		onSelection(buildables[selection])
	})
}


func (p *Plot) Tick(tick int) {
	p.Unit.Resources.Current += p.Unit.Resources.Rate
	p.Updated = tick
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


func (w *World) ChangedPlots(plots []*Plot, tick int, playerMsk int64) []*Plot {
	for x := 0; x < len(w.Plots); x += 1 {
		for y := 0; y < len(w.Plots[x]); y += 1 {
			plot := &w.Plots[x][y]

			if plot.Explored | playerMsk > 0 && plot.Updated == tick {
				plots = append(plots, plot)
			}
		}
	}

	return plots
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


func (w *World) Reveal(x, y, r int, playerId int64) Region {
	min_x, max_x := max(x - r * 2, 0), min(x + r * 2, len(w.Plots)-1)
	min_y, max_y := max(y - r, 0), min(y + r, len(w.Plots[0])-1)

	for i := min_x; i <= max_x; i += 1 {
		for j := min_y; j <= max_y; j += 1 {
			dx, dy := (x - i) / 2, y - j
			// dist := math.Sqrt(math.Pow(float64(x-i) / 2, 2) + math.Pow(float64(y-j), 2))
			dist := dx * dx + dy * dy
			if dist <= (r * r) {
				w.Plots[i][j].Explored |= playerId
			}

			// w.Plots[i][j].Explored = 0
			// w.Plots[i][j].Explored |= playerId
		}
	}

	// min_x, min_y = 0, 0
	// max_x, max_y = 199, 199

	return Region{ Min: Point{ min_x, min_y }, Max: Point{ max_x, max_y }}
}


func (w *World) Tick(tick int) {
	width, height := len(w.Plots), len(w.Plots[0])
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			w.Plots[x][y].Tick(tick)
		}
	}
}


func (w *World) Init(seed int64) {
	if w.R == nil {
		w.R = rand.New(rand.NewSource(seed))
	}

	// populate with initial noise
	width, height := len(w.Plots), len(w.Plots[0])
	unitNone := UnitIndex("vacant")

	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			plot := &w.Plots[x][y]

			plot.Unit = Units[unitNone]
			plot.X, plot.Y = x, y
			plot.Explored = 0
			plot.Unit.OwnerID = 0

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
