package main

import (
	"github.com/nsf/termbox-go"
)

type Icon struct {
	Symbol rune
	Attr termbox.Attribute
}


type ResourcesProps struct {
	Current float32
	Rate float32
}


type PlotUnit struct {
	Icon
	Resources ResourcesProps
	Owner *Player

}
type PlotTile struct { Icon }

const (
	UnitNone = 0
	UnitForest
	UnitVillage
	UnitCity
	UnitFarm
	UnitMine
	UnitExplorers
	UnitTravelers
	UnitMerchants
	UnitArmy
)

var Units = []PlotUnit{
	PlotUnit{Icon{rune(' '), 0}, ResourcesProps{Current:0, Rate:0}, nil},
	PlotUnit{Icon{rune('^'), termbox.AttrUnderline}, ResourcesProps{Current:0, Rate:1}, nil},
	PlotUnit{Icon{rune('∆'), termbox.AttrUnderline}, ResourcesProps{Current:0, Rate:3}, nil},
	PlotUnit{Icon{rune('U'), 0}, ResourcesProps{Current:0, Rate:10}, nil},
	PlotUnit{Icon{rune('≈'), 0}, ResourcesProps{Current:0, Rate:2}, nil},
	PlotUnit{Icon{rune('M'), 0}, ResourcesProps{Current:0, Rate:3}, nil},
	PlotUnit{Icon{rune('*'), 0}, ResourcesProps{Current:0, Rate:-1}, nil},
	PlotUnit{Icon{rune('*'), 0}, ResourcesProps{Current:0, Rate:0}, nil},
	PlotUnit{Icon{rune('*'), termbox.AttrBold}, ResourcesProps{Current:0, Rate:5}, nil},
	PlotUnit{Icon{rune('*'), termbox.AttrBold}, ResourcesProps{Current:0, Rate:-5}, nil},
}


const (
	TileNone = 0
	TileRoad
)

var Tiles = []PlotTile{
	PlotTile{Icon{rune(' '), 0}},
	PlotTile{Icon{rune(' '), 0}},
}
