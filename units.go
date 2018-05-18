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
	Type int
	Resources ResourcesProps
	Owner *Player
	Name string
}
type PlotTile struct { Icon }

const (
	UnitNone      = 0
	UnitForest    = 1
	UnitVillage   = 2
	UnitCity      = 3
	UnitFarm      = 4
	UnitMine      = 5
	UnitExplorers = 6
	UnitTravelers = 7
	UnitMerchants = 8
	UnitArmy      = 9
)

var Units = []PlotUnit{
	PlotUnit{Icon{rune(' '), 0}, UnitNone, ResourcesProps{Current:0, Rate:0}, nil, "vacant"},
	PlotUnit{Icon{rune('^'), termbox.AttrUnderline}, UnitForest, ResourcesProps{Current:0, Rate:1}, nil, "forest"},
	PlotUnit{Icon{rune('∆'), termbox.AttrUnderline}, UnitVillage, ResourcesProps{Current:0, Rate:3}, nil, "village"},
	PlotUnit{Icon{rune('U'), 0}, UnitCity, ResourcesProps{Current:0, Rate:10}, nil, "city"},
	PlotUnit{Icon{rune('≈'), 0}, UnitFarm, ResourcesProps{Current:0, Rate:2}, nil, "farm"},
	PlotUnit{Icon{rune('M'), 0}, UnitMine, ResourcesProps{Current:0, Rate:3}, nil, "mine"},
	PlotUnit{Icon{rune('*'), 0}, UnitExplorers, ResourcesProps{Current:0, Rate:-1}, nil, "explorers"},
	PlotUnit{Icon{rune('*'), 0}, UnitTravelers, ResourcesProps{Current:0, Rate:0}, nil, "travelers"},
	PlotUnit{Icon{rune('*'), termbox.AttrBold}, UnitMerchants, ResourcesProps{Current:0, Rate:5}, nil, "merchants"},
	PlotUnit{Icon{rune('*'), termbox.AttrBold}, UnitArmy, ResourcesProps{Current:0, Rate:-5}, nil, "army"},
}


const (
	TileNone = 0
	TileRoad
)

var Tiles = []PlotTile{
	PlotTile{Icon{rune(' '), 0}},
	PlotTile{Icon{rune(' '), 0}},
}

const (
	UnitDescriptionShort = 0
	UnitDescriptionFull
)

func (u *PlotUnit) Description(descType int) string {
	desc := u.Name

	if u.Owner != nil {
		desc = u.Owner.Name + " " + desc
	}

	return desc
}
