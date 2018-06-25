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
	IdxSea      = 0
	IdxBeach    = 1
	IdxPlains   = 2
	IdxForrest  = 3
	IdxMountain = 4
)

var PlotTypes = [...]float32{ PlotSea, PlotBeach, PlotPlains, PlotForest, PlotMountain }
var PlotTypeNames = [...]string{ "sea", "beach", "plains", "forrest", "mountains" }

const (
	MaxWorldWidth = 100
	MaxWorldHeight = 50
	ViewWidth = 80
	ViewHeight = 32
)

type Plot struct {
	Explored int64
	Elevation float32
	Productivity float32
	OwnerID int64
	Unit PlotUnit
	Tile PlotTile
	X, Y int
	Updated int
	Tag int16
}

type World struct {
	Width, Height int
	Plots [MaxWorldWidth][MaxWorldHeight]Plot
	Smoothness int
	R *rand.Rand
}


func TerrainIndexForElevation(elevation float32) int {
	for i := 0; i < len(PlotTypes); i += 1 {
		if elevation <= PlotTypes[i] {
			return i
		}
	}

	return -1
}


func (p *Plot) TerrainIndex() int {
	return TerrainIndexForElevation(p.Elevation)
}


func (p *Plot) TerrainName() string {
	return PlotTypeNames[p.TerrainIndex()]
}


func (p *Plot) Neighbors(world *World) []*Plot {
	slice := make([]*Plot, 0, 8)

	for i := -1; i <= 1; i += 1 {
		for j := -1; j <= 1; j += 1 {
			r, c := p.X + i, p.Y + j
			if i | j == 0 { continue }
			if r < 0 || r >= MaxWorldWidth { continue }
			if c < 0 || c >= MaxWorldHeight { continue }
			n := &world.Plots[r][c]
			// GfxMsg(fmt.Sprintf("(%d, %d) -> (%d, %d) unit -> %d", r, c, n.X, n.Y, n.Unit.Type))

			slice = append(slice, n)
		}
	}

	return slice
}


func (p *Plot) IsBorderBoundry(w *World) (playerId int64, present int) {
	playerId = p.OwnerID

	if playerId == 0 { return 0, 0 }

	neighbors := p.Neighbors(w)
	borderCount := 0

	for ni := 0; ni < len(neighbors); ni += 1 {
		if neighbors[ni].OwnerID == playerId {
			borderCount += 1
		}
	}

	if borderCount < 8 {
		present = 1
	} else {
		present = 0
	}

	return
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


func (p *Plot) HasNeighboringPlotType(world *World, neighbors []*Plot, plotTypeIndex int) bool {
	for ni := 0; ni < len(neighbors); ni += 1 {
		neighbor := neighbors[ni]
		if neighbor.TerrainIndex() == plotTypeIndex {
			return true
		}
	}

	return false
}


func (p *Plot) PossibleBuilds(world *World, owner int64) []int {
	unitCity := UnitIndex("city")
	unitVillage := UnitIndex("village")
	unitFarm := UnitIndex("farm")
	unitRoad := UnitIndex("road")
	unitBridge := UnitIndex("bridge")
	unitMine := UnitIndex("mine")
	unitCanal := UnitIndex("canal")

	neighbors := p.Neighbors(world)
	nextToVillage := p.HasNeighbor(world, neighbors, unitVillage, owner)
	nextToCity := p.HasNeighbor(world, neighbors, unitCity, owner)
	nextToFarm := p.HasNeighbor(world, neighbors, unitFarm, owner)
	nextToRoad := p.HasNeighbor(world, neighbors, unitRoad, owner)

	nextToWater := p.HasNeighboringPlotType(world, neighbors, IdxSea)

	buildables := make([]int, 0, 10)

	nextToCivilization := nextToRoad || nextToCity || nextToVillage || nextToFarm

	// TODO: figure out why this causes problems on the
	// client, but not the server
	if p.Unit.Type != UnitIndex("vacant") {
		return buildables
	}

	if nextToWater && p.TerrainIndex() != IdxSea {
		if nextToCity || nextToVillage || nextToRoad {
			buildables = append(buildables, unitCity)
			buildables = append(buildables, unitVillage)
			buildables = append(buildables, unitCanal)
		}

		if nextToCivilization && p.TerrainIndex() == IdxPlains {
			buildables = append(buildables, unitFarm)
		}
	}

	if nextToCivilization {
		if p.TerrainIndex() == IdxSea {
			buildables = append(buildables, unitBridge)
		} else {
			buildables = append(buildables, unitRoad)
		}

		if p.TerrainIndex() == TerrainIndexForElevation(PlotMountain) {
			buildables = append(buildables, unitMine)
		}
	}

	return buildables
}


func (p *Plot) Description(descType int) string {
	desc := fmt.Sprintf("%v x%0.2f (%d,%d)", p.TerrainName(), p.Productivity, p.X, p.Y)
	landDesc := " in "
	if p.OwnerID > 0 {
		landDesc += PlayerFromID(p.OwnerID).Name + "'s "
	} else {
		landDesc += " the "
	}
	landDesc += desc

	if p.Unit.Type != UnitIndex("vacant") {
		return p.Unit.Description(descType) + landDesc
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

	if p.Unit.Type == UnitIndex("canal") {
		p.Elevation = PlotTypes[IdxSea] 
	}

	fmt.Println(msg)

	return &p.Unit, msg
}


func (p *Plot) BuildMenu(world *World, owner int64, onSelection func(int)) {
	buildables := p.PossibleBuilds(world, owner)
	options := [10]string{}
	optSlice := options[0:0]

	for i := 0; i < len(buildables); i += 1 {
		ui := buildables[i]
		optSlice = append(optSlice, fmt.Sprintf("%s ($%0.2f)", Units[ui].Name, Units[ui].Resources.Cost))
	}

	GfxMenu("These units can be built here", optSlice, func(selection int) {
		onSelection(buildables[selection])
	})
}


func (p *Plot) AvailableResources(world *World, playerId int64, tag int16) float32 {
	available := float32(0.0)

	if p.Tag != tag && p.Unit.OwnerID == playerId {
		p.Tag = tag
		available += p.Unit.Resources.Current
	}

	neighbors := p.Neighbors(world)
	for i := 0; i < len(neighbors); i += 1 {
		neighbor := neighbors[i]
		if neighbor.Tag != tag && neighbor.Unit.OwnerID == playerId {
			available += neighbor.AvailableResources(world, playerId, tag)
		}
	}

	return available
}


func (p *Plot) DeductResources(world *World, amount float32, playerId int64, tag int16) {
	if p.Tag != tag && p.Unit.OwnerID == playerId {
		p.Tag = tag

		if amount > p.Unit.Resources.Current {
			amount -= p.Unit.Resources.Current
			p.Unit.Resources.Current = 0
		} else {
			p.Unit.Resources.Current -= amount
			return
		}
	}

	neighbors := p.Neighbors(world)
	for i := 0; i < len(neighbors); i += 1 {
		neighbor := neighbors[i]
		if neighbor.Tag != tag && neighbor.Unit.OwnerID == playerId {
			neighbor.DeductResources(world, amount, playerId, tag)
		}
	}
}

func (p *Plot) SpendResources(world *World, amount float32, playerId int64) bool {
	available := p.AvailableResources(world, playerId, int16(rand.Int31()))

	if available >= amount {
		p.DeductResources(world, amount, playerId, int16(rand.Int31()))
		return true
	}

	return false
}


func (p *Plot) Tick(tick int) {
	unitNone := UnitIndex("vacant")
	if p.Unit.Type != unitNone {
		p.Unit.Resources.Current += p.Unit.Resources.Rate * p.Productivity

		if p.Unit.Resources.Current < 0 {
			p.Unit.Type = unitNone
		}
	} else {
		p.Unit.Resources.Current = 0
	}

	p.Updated = tick
}


func (p *Plot) PlayerResources(playerId int64) (current, rate float32) {
	if p.Unit.OwnerID != playerId { return 0, 0 }

	return p.Unit.Resources.Current, p.Unit.Resources.Rate * p.Productivity
}



func (p *Plot) ProductionRate() float32 {
	switch {
	case p.Elevation <= PlotSea:
		return float32(1.25 / 9.0)
	case p.Elevation <= PlotBeach:
		return float32(0.0 / 9.0)
	case p.Elevation <= PlotPlains:
		return float32(1.0 / 9.0)
	case p.Elevation <= PlotForest:
		return float32(1.25 / 9.0)
	case p.Elevation <= PlotMountain:
		return float32(0.75 / 9.0)
	}

	return 0;
}


func (p *Plot) IsLivable() bool {
	return p.Elevation > PlotSea && p.Elevation < PlotForest
}


func (w *World) ChangedPlots(plots []*Plot, tick int, playerMsk int64) []*Plot {
	for x := 0; x < w.Width; x += 1 {
		for y := 0; y < w.Height; y += 1 {
			plot := &w.Plots[x][y]

			if plot.Explored | playerMsk > 0 && plot.Updated == tick {
				plots = append(plots, plot)
			}
		}
	}

	return plots
}


func (w *World) TotalCaptureSpace() (captureable int) {
	captureable = w.Width * w.Height

	return
}


func (w *World) IsGameOver() bool {
	for x := 0; x < w.Width; x += 1 {
		for y := 0; y < w.Height; y += 1 {
			plot := &w.Plots[x][y]

			if plot.OwnerID == 0 {
				return false
			}
		}
	}

	return true
}


func (w *World) PlayerScore(playerId int64) int {
	score := 0
	for x := 0; x < w.Width; x += 1 {
		for y := 0; y < w.Height; y += 1 {
			plot := &w.Plots[x][y]

			if plot.OwnerID == playerId {
				score += 1
			}
		}
	}

	return score
}


func (w *World) FindLivablePlot() *Plot {
	for {
		x, y := w.R.Intn(MaxWorldWidth), w.R.Intn(MaxWorldHeight)
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
	min_x, max_x := max(x - r * 2, 0), min(x + r * 2, w.Width-1)
	min_y, max_y := max(y - r, 0), min(y + r, len(w.Plots[0])-1)

	for i := min_x; i <= max_x; i += 1 {
		for j := min_y; j <= max_y; j += 1 {
			dx, dy := (x - i) / 2, y - j
			// dist := math.Sqrt(math.Pow(float64(x-i) / 2, 2) + math.Pow(float64(y-j), 2))
			dist := dx * dx + dy * dy
			if dist <= (r * r) {
				plot := &w.Plots[i][j]
				if plot.OwnerID == 0 {
					plot.OwnerID = playerId
				}

				plot.Explored |= playerId
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
	width, height := w.Width, len(w.Plots[0])
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			w.Plots[x][y].Tick(tick)
		}
	}
}


func (w *World) PlayerResources(playerId int64) (current, rate float32) {

	width, height := w.Width, len(w.Plots[0])
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			c, r := w.Plots[x][y].PlayerResources(playerId)
			current += c
			rate += r
		}
	}

	return
}


func (w *World) Init(seed int64) {
	if w.R == nil {
		w.R = rand.New(rand.NewSource(seed))
	}

	// populate with initial noise
	width, height := w.Width, w.Height
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
		for x := 1; x < w.Width - 1; x += 1 {
			for y := 1; y < w.Height - 1; y += 1 {
				plot := &w.Plots[x][y]
				plot.Elevation = w.avgPatchElevation(x, y)
				// fmt.Println(patch)
			}
		}
	}

	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			plot := &w.Plots[x][y]
			neighbors := plot.Neighbors(w)

			plot.Productivity = plot.ProductionRate()
			for i := 0; i < len(neighbors); i += 1 {
				plot.Productivity += neighbors[i].ProductionRate();
			}
		}
	}
}
