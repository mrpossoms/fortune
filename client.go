package main


import (
	"github.com/nsf/termbox-go"
	"net"
	"os"
	"sync"
	// "fmt"
)


func GameClient() {
	var joinLock = &sync.Mutex{}
	var player Player;
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		// handle error
		panic(err)
	}

	joinLock.Lock()

	go func(){
		for {
			hdr := Msg{}
			hdr.Read(conn)

			switch (hdr.Type) {
			case PayTypJoin:
				player.Read(conn)

				GfxDrawBegin()
				player.Name = GfxPrompt("Type your name")

				hdr.Write(conn)
				player.Write(conn)
				joinLock.Unlock()
				break
			case PayTypPlot:
				var plot Plot
				for i := 0; i < int(hdr.Count); i += 1 {
					plot.Read(conn)
					GameWorld.Plots[plot.X][plot.Y] = plot
				}
				break
			}
		}

		conn.Close()
	}()

	joinLock.Lock()
	running := true

	GfxInit()
	for running {
		GfxDrawBegin()

		GameWorld.GfxDraw(&player)

		evt := GfxDrawFinish(true)

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
			GameWorld.Plots[x][y].Explored = 1
		}
	}
/*
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

	start:=GameWorld.FindLivablePlot()

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

		GameWorld.GfxDraw(&player)

		evt := GfxDrawFinish(true)

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
			GameWorld.Plots[x][y].Explored = 1
		}
	}

	GfxUninit()
*/
}
