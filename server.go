package main

import (
	"time"
	"net"
	"fmt"
	"encoding/gob"
)


func GameServer(ln net.Listener) {
	var playerPool [64]Player

	players := playerPool[0:0]

	GameWorld = World {
		Smoothness: 2,
	}

	//Generate the landscape
	fmt.Print("Generating map...")
	for i := 0; i < 3; i += 1 {
		GameWorld.Init(int64(time.Now().Second()))
	}
	fmt.Println("DONE")

	// Accept and handle connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}

		fmt.Println("Connection incoming")
		go func() {
			msg := Msg{ Type: PayTypJoin, Count:1 }
			player := Player{ ID: 1 << uint(len(players))}

			pconn := PlayerConnection{
				Conn: conn,
				Index: PlayerIndex(player.ID),
				ID: player.ID,
				Enc: gob.NewEncoder(conn),
				Dec: gob.NewDecoder(conn),
			}
			PlayerConns[PlayerIndex(player.ID)] = pconn

			// Find a place for them to start
			start:=GameWorld.FindLivablePlot()

			// player view and cam setup
			player.Cam.X, player.Cam.Y = WorldWidth / 2, WorldHeight / 2
			player.Cursor.X, player.Cursor.Y = player.Cam.X, player.Cam.Y
			player.Cam.View.Width = ViewWidth
			player.Cam.View.Height = ViewHeight
			player.MoveCursorTo(start.X, start.Y)

			// Send the initial empty player object
			if err := msg.Write(pconn.Enc); err != nil { panic(err) }
			if err := player.Write(pconn.Enc); err != nil { panic(err) }

			// Continuous message handling
			for {
				msg.Read(pconn.Dec)

				switch {
				case msg.Type == PayTypJoin:
					player.Read(pconn.Dec)
					players = append(players, player)
					fmt.Println(players[len(players)-1].Name + " has joined the game")

					min_x, min_y, max_x, max_y := GameWorld.Reveal(start.X, start.Y, 4, player.ID)
					fmt.Printf("(%d, %d) -> (%d, %d)\n", min_x, min_y, max_x, max_y)

					// Send only their visible part of the map
					Msg{ Type: PayTypPlot, Count: int32((1 + max_x - min_x) * (1 + max_y - min_y)) }.Write(pconn.Enc)
					for x := min_x; x <= max_x; x += 1 {
						for y := min_y; y <= max_y; y += 1 {
							plot:=&GameWorld.Plots[x][y]

							if err := plot.Write(pconn.Enc); err != nil {
								fmt.Println("Sending map failed!")
							}
						}
					}

					// Spawn their village
					start.SpawnUnit(UnitVillage, &player)
					Msg{ Type: PayTypPlot, Count: 1 }.Write(pconn.Enc)
					start.Write(pconn.Enc)

					// for j := 0; j < len(players); j += 1 {
					// 	Msg{ Type: PayTypPlayer, Count: int32(len(players))}.Write(PlayerConns[j].Enc)
					// 	for i := 0; i < len(players); i += 1 {
					// 		players[i].Write(PlayerConns[j].Enc)
					// 	}
					// }


					// Send all the players
					// Msg{ Type: PayTypPlayer, Count: int32(len(players))}.Broadcast(players, func (c *PlayerConnection) {
					// 	for i := 0; i < len(players); i += 1 {
					// 		players[i].Write(c.Enc)
					// 	}
					// })
					//

					Msg{ Type: PayTypPlayer, Count: int32(len(players))}.Broadcast(players)
					for i := 0; i < len(players); i += 1 {
						players[i].Broadcast(players)
					}

				case msg.Type == PayTypPlot:
					updatedPlot := Plot{}
					updatedPlot.Read(pconn.Dec)
					x, y := updatedPlot.X, updatedPlot.Y

					fmt.Printf("Got plot (%d,%d)", x, y)

					GameWorld.Plots[x][y] = updatedPlot
					Msg { Type: PayTypPlot, Count: 1 }.Write(pconn.Enc)
					updatedPlot.Write(pconn.Enc)

					break
				}

				msg.Type = -1
				// fmt.Println("Got something")
			}

			conn.Close()
		}()
	}
}
