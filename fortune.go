package main

import (
	// "fmt"
	"time"
	"net"
	"github.com/nsf/termbox-go"
)


func main() {
	ln, err := net.Listen("tcp", ":31337")

	world := World {
		Smoothness: 2,
	}

	if err != nil {
		// Can't host, act as a client

		panic(err)
	} else {
		ln.Close()
		// Generate the landscape
		// for i := 0; i < 3; i += 1 {
		// 	world.Init(int64(time.Now().Second()))
		// }
		//
		// // Accept and handle connections
		// for {
		//
		// }
	}

	for i := 0; i < 3; i += 1 {
		world.Init(int64(time.Now().Second()))
	}

	player := Player {
		Name: "mrpossoms",
		ID: 0x1,
		Colors: PlayerColors{
			Fg: termbox.ColorBlack,
			Bg: termbox.ColorRed,
		},
	}

	player.Cam.X, player.Cam.Y = WorldWidth / 2, WorldHeight / 2
	player.Cam.View.Width = ViewWidth
	player.Cam.View.Height = ViewHeight


	start:=world.FindLivablePlot()

	unit, msg := start.SpawnUnit(UnitVillage, &player)

	if unit != nil {
		player.MoveCursorTo(start.X, start.Y)
	}
	GfxMsg(msg)


	GfxInit()
	GfxMsg("Hello, world")



	running := true

	for running {
		GfxDrawBegin()

		world.GfxDraw(&player)

		evt := GfxDrawFinish()

		switch evt.Key {
		case termbox.KeyArrowUp:
			player.MoveCursor(0, -1)
			break
		case termbox.KeyArrowDown:
			player.MoveCursor(0, 1)
			break
		case termbox.KeyArrowLeft:
			player.MoveCursor(-1, 0)
			break
		case termbox.KeyArrowRight:
			player.MoveCursor(1, 0)
			break
		}

		switch evt.Ch {
		case rune('q'):
			running = false
			break
		}

		{
			x, y := player.Cursor.X, player.Cursor.Y
			world.Plots[x][y].Explored = 1
		}
	}

	GfxUninit()
}
