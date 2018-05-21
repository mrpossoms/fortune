package main


import (
	"github.com/nsf/termbox-go"
	"net"
	"os"
	"golang.org/x/sync/semaphore"
	"context"
	// "fmt"
	"encoding/gob"
)


func GameClient() {

	if len(os.Args) < 2 {
		panic("Please provide an address or domain to connect to")
	}

	ctx := context.TODO()
	joinSem := semaphore.NewWeighted(2)
	var player Player;
	gotMap := false
	conn, err := net.Dial("tcp", os.Args[1] + ":31337")
	if err != nil {
		// handle error
		panic(err)
	}

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	joinSem.Acquire(ctx, 2)

	go func(){
		for {
			hdr := Msg{}
			if err:= hdr.Read(dec); err != nil {
				panic(err)
			}

			switch (hdr.Type) {
			case PayTypJoin:
				player.Read(dec)

				GfxInit()
				GfxDrawBegin()
				// player.Name = "mrpossoms"//GfxPrompt("Type your name")
				player.Name = GfxPrompt("Type your name")


				hdr.Write(enc)
				player.Write(enc)
				joinSem.Release(1)
				break
			case PayTypPlot:
				var plot Plot
				for i := 0; i < int(hdr.Count); i += 1 {
					plot.Read(dec)
					GameWorld.Plots[plot.X][plot.Y] = plot
				}

				if !gotMap {
					gotMap = true
					joinSem.Release(1)
				} else {
					termbox.Interrupt()
				}

				break
			case PayTypPlayer:
				player := Player{}
				for i := 0; i < int(hdr.Count); i += 1 {
					player.Read(dec)
					*PlayerFromID(player.ID) = player
				}
				break
			}

			termbox.Interrupt()
		}

		conn.Close()
	}()

	joinSem.Acquire(ctx, 1)
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
