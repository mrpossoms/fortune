package main

import (
	"github.com/nsf/termbox-go"
)


type Box struct {
	X, Y int
	View struct {
		Width, Height int
	}
}

type Camera struct {
	Box
}


type PlayerColors struct {
	Bg termbox.Attribute
	Fg termbox.Attribute
}


type Player struct {
	Score int32
	ID int64
	Income float32
	Wealth float32
	Colors PlayerColors
	Cursor Box
	Cam Camera
	Name string
}


func PlayerFromID(id int64) *Player {
	return &Players[PlayerIndex(id)]
}


func PlayerIndex(id int64) int {
	i := 0

	for id >>= 1; id > 0; i += 1 {
		id >>= 1
	}

	return i
}


func (p *Player) MoveCursor(dx, dy int) {
	last_dist_x := abs(p.Cursor.X - p.Cam.X)
	last_dist_y := abs(p.Cursor.Y - p.Cam.Y)

	for i := abs(dy); i > 0; i -= 1 {
		_dy := dy / abs(dy)
		p.Cursor.Move(0, _dy)
		new_dist := abs(p.Cursor.Y - p.Cam.Y)
		if new_dist > 10 && new_dist > last_dist_y {
			p.Cam.Move(0, _dy)
		}
	}

	for i := abs(dx); i > 0; i -= 1 {
		_dx := dx / abs(dx)
		p.Cursor.Move(_dx, 0)
		new_dist := abs(p.Cursor.X - p.Cam.X)
		if new_dist > 10 && new_dist > last_dist_x {
			p.Cam.Move(_dx, 0)
		}
	}
}


func (p *Player) MoveCursorTo(x, y int) {
	p.MoveCursor(x - p.Cursor.X, y - p.Cursor.Y)
}


func (p *Player) SelectedPlot(world *World) *Plot {
	return &world.Plots[p.Cursor.X][p.Cursor.Y]
}


func (c *Box) Move(dx, dy int) {
	nx, ny := c.X + dx, c.Y + dy

	if nx >= 0 && nx < GameWorld.Width {
		c.X = nx
	}

	if ny >= 0 && ny < GameWorld.Height {
		c.Y = ny
	}
}


func (c *Camera) Move(dx, dy int) {
	hw, hh := c.View.Width >> 1, c.View.Height >> 1
	nx, ny := c.X + dx, c.Y + dy

	if nx - hw >= 0 && nx + hw < GameWorld.Width {
		c.X = nx
	}

	if ny - hh >= 0 && ny + hh < GameWorld.Height {
		c.Y = ny
	}
}
