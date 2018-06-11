package main

import (
	"github.com/nsf/termbox-go"
	// "fmt"
)

type Icon struct {
	Symbol rune
	Attr termbox.Attribute
}


type ResourcesProps struct {
	Current float32
	Rate float32
	Cost float32
}


type PlotUnit struct {
	Icon
	Type int
	Resources ResourcesProps
	OwnerID int64
	Name string
}
type PlotTile struct { Icon }

// const (
// 	UnitNone      = 0
// 	UnitFarm      = 1
// 	UnitMine      = 2
// 	UnitForest    = 3
// 	UnitCity      = 4
// 	UnitVillage   = 5
// 	UnitExplorers = 6
// 	UnitTravelers = 7
// 	UnitMerchants = 8
// 	UnitArmy      = 9
// )

var Units = []PlotUnit{
	PlotUnit{Icon{rune(' '), 0}, 0, ResourcesProps{Current:0, Rate:0}, 0, "vacant"},
	PlotUnit{Icon{rune('≈'), 0}, 0, ResourcesProps{Current:0, Rate:2, Cost:30}, 0, "farm"},
	PlotUnit{Icon{rune('M'), 0}, 0, ResourcesProps{Current:0, Rate:3, Cost:400}, 0, "mine"},
	PlotUnit{Icon{rune('^'), termbox.AttrUnderline}, 0, ResourcesProps{Current:0, Rate:1}, 0, "forest"},
	PlotUnit{Icon{rune('Ü'), 0}, 0, ResourcesProps{Current:0, Rate:10, Cost: 3000}, 0, "city"},
	PlotUnit{Icon{rune('∆'), termbox.AttrUnderline}, 0, ResourcesProps{Current:0, Rate:1, Cost: 300}, 0, "village"},
	PlotUnit{Icon{rune('*'), 0}, 0, ResourcesProps{Current:0, Rate:-1}, 0, "explorers"},
	PlotUnit{Icon{rune('*'), 0}, 0, ResourcesProps{Current:0, Rate:0}, 0, "travelers"},
	PlotUnit{Icon{rune('*'), termbox.AttrBold}, 0, ResourcesProps{Current:0, Rate:5}, 0, "merchants"},
	PlotUnit{Icon{rune('*'), termbox.AttrBold}, 0, ResourcesProps{Current:0, Rate:-5}, 0, "army"},
	PlotUnit{Icon{rune('+'), termbox.AttrBold}, 0, ResourcesProps{Current:0, Rate:0, Cost: 10}, 0, "road"},
}


var unitMap map[string]int


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


func UnitIndex(name string) int {
	if unitMap == nil {
		unitMap = make(map[string]int)

		for i := 0; i < len(Units); i += 1 {
			unitMap[Units[i].Name] = i
			Units[i].Type = i
		}
	}

	return unitMap[name]
}


func (u *PlotUnit) Description(descType int) string {
	desc := u.Name

	owner := PlayerFromID(u.OwnerID)
	if owner != nil {
	 	desc = owner.Name + "'s " + desc
	}

	return desc
}
