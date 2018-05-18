package main

import (
	"github.com/nsf/termbox-go"
)

type Player struct {
	Name string
	Score int
	ID int32
	Color termbox.Attribute
}
