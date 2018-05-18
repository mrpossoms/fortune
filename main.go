package main

import (
	// "fmt"
	"time"
	"github.com/nsf/termbox-go"
)

func abs(i int) int {
	if i < 0 { return -i }
	return i
}

func main() {
	world := World {
		Smoothness: 2,
	}

	for i := 0; i < 3; i += 1 {
		world.Init(int64(time.Now().Second()))
	}


	GfxInit()
	GfxMsg("Hello, world")

	running := true

	for running {
		GfxDrawBegin()

		world.GfxDraw()

		evt := GfxDrawFinish()

		last_dist_x := abs(world.Cursor.X - world.Cam.X)
		last_dist_y := abs(world.Cursor.Y - world.Cam.Y)

		switch evt.Key {
		case termbox.KeyArrowUp:
			world.Cursor.Move(0, -1)
			new_dist := abs(world.Cursor.Y - world.Cam.Y)
			if new_dist > 10 && new_dist > last_dist_y {
				world.Cam.Move(0, -1)
			}
			break
		case termbox.KeyArrowDown:
			world.Cursor.Move(0, 1)
			new_dist := abs(world.Cursor.Y - world.Cam.Y)
			if new_dist > 10 && new_dist > last_dist_y {
				world.Cam.Move(0, 1)
			}
			break
		case termbox.KeyArrowLeft:
			world.Cursor.Move(-1, 0)
			new_dist := abs(world.Cursor.X - world.Cam.X)
			if new_dist > 25 && new_dist > last_dist_x {
				world.Cam.Move(-1, 0)
			}
			break
		case termbox.KeyArrowRight:
			world.Cursor.Move(1, 0)
			new_dist := abs(world.Cursor.X - world.Cam.X)
			if new_dist > 25 && new_dist > last_dist_x {
				world.Cam.Move(1, 0)
			}
			break
		}

		switch evt.Ch {
		case rune('q'):
			running = false
			break
		}

		{
			x, y := world.Cursor.X, world.Cursor.Y
			world.Plots[x][y].Explored = 1
		}
	}

	GfxUninit()
}
