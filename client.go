package main


import (
	"github.com/nsf/termbox-go"
	"net"
	"os"
	"golang.org/x/sync/semaphore"
	"context"
	"fmt"
	"encoding/gob"
	// "time"
)


func GameClient() {

	if len(os.Args) < 2 {
		panic("Please provide an address or domain to connect to")
	}

	ctx := context.TODO()
	joinSem := semaphore.NewWeighted(2)
	var player *Player;
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
		msg := Msg{}
		for {
			// TODO: figue out why this is getting out of sync and
			// not reading headers when building near a tick
			if err:= msg.Read(dec); err != nil {
				panic(err)
			}

			switch (msg.Type) {
			case PayTypJoin:
				p := Player{}
				p.Read(dec)

				*PlayerFromID(p.ID) = p
				player = PlayerFromID(p.ID)

				GfxInit()
				GfxDrawBegin()
				player.Name = GfxPrompt("Type your name")

				colorOptions := []string{ "red", "blue", "black", "magenta", "cyan" }
				colors := []PlayerColors{
					PlayerColors{ termbox.ColorRed, termbox.ColorWhite },
					PlayerColors{ termbox.ColorBlue, termbox.ColorWhite },
					PlayerColors{ termbox.ColorBlack, termbox.ColorWhite },
					PlayerColors{ termbox.ColorMagenta, termbox.ColorWhite },
					PlayerColors{ termbox.ColorCyan, termbox.ColorWhite },
				}

				waiting := true

				GfxMenu("Choose your color", colorOptions, func(selected int) {
					waiting = false
					player.Colors = colors[selected]
				})

				GfxDrawBegin()
				for ; waiting ; {
					GfxDrawFinish(true)
				}

				msg.Write(enc)
				player.Write(enc)
				joinSem.Release(1)

				break
			case PayTypInfo:
				info := GameInfo{}
				info.Read(dec)
				GameWorld.Width = info.Width
				GameWorld.Height = info.Height
				fmt.Printf("got info %dx%d\n", info.Width, info.Height)
				break
			case PayTypPlot:
				for i := 0; i < int(msg.Count); i += 1 {
					plot := Plot{}
					plot.Read(dec)
					GameWorld.Plots[plot.X][plot.Y] = plot
				}

				if !gotMap {
					gotMap = true
					joinSem.Release(1)
				}

				break
			case PayTypPlayer:
				p := Player{}
				for i := 0; i < int(msg.Count); i += 1 {
					p.Read(dec)

					localPlayer := PlayerFromID(p.ID)

					if localPlayer.Name != p.Name {
						GfxMsg(fmt.Sprintf("%v joined the game", p.Name))
					}

					if p.ID == localPlayer.ID {
						localPlayer.Colors = p.Colors
						localPlayer.Wealth = p.Wealth
						localPlayer.Income = p.Income
						localPlayer.Score = p.Score
					} else {
						*PlayerFromID(p.ID) = p
					}


					// fmt.Println(player.Name + " joined the game")
				}
				break
			case PayTypText:
				text := TextPayload{}
				text.Read(dec)
				GfxMsg(text.Msg)
				break
			}

			if player != nil {
				termbox.Interrupt()
			}

			msg.Type = -1

		}

		conn.Close()
	}()

	joinSem.Acquire(ctx, 1)
	running := true

	GfxInit()

	for running {
		GfxDrawBegin()

		GameWorld.GfxDraw(player)

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
		case termbox.KeyTab:
			_, h := termbox.Size()
			score_board := ""
			for i := 0; i < len(Players); i += 1 {
				if Players[i].ID > 0 {
					score_board += fmt.Sprintf("%v - score %d\n", Players[i].Name, Players[i].Score)
				}
			}

			GfxMsgExplicit(MsgContainer { Str: score_board, Y: h / 2 })
		}

		x, y := player.Cursor.X, player.Cursor.Y
		selectedPlot := &GameWorld.Plots[x][y]

		switch evt.Ch {
		case rune('q'):
			running = false
			break
		case rune('b'):
			selectedPlot.BuildMenu(&GameWorld, player.ID, func(selection int) {
				updatedPlot := *selectedPlot
				updatedPlot.Unit = Units[selection]
				updatedPlot.Unit.OwnerID = player.ID
				Msg { Type: PayTypPlot, Count: 1 }.Write(enc)
				updatedPlot.Write(enc)
			})
		}

		{
			// x, y := player.Cursor.X, player.Cursor.Y
			// GameWorld.Plots[x][y].Explored = 1
		}
	}
}
