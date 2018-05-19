package main

import (
	"time"
	"net"
	"fmt"
	"encoding/gob"
)


func GameServer(ln net.Listener) {
	var playerPool [32]Player
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

		// connHandler :=
		enc := gob.NewEncoder(conn)
		dec := gob.NewDecoder(conn)

		fmt.Println("Connection incoming")
		go func() {
			msg := Msg{ Type: PayTypJoin, Count:1 }
			player := Player{ ID: 1 << uint(len(players))}

			// Find a place for them to start
			start:=GameWorld.FindLivablePlot()

			min_x, min_y, max_x, max_y := GameWorld.Reveal(start.X, start.Y, 4, player.ID)
			fmt.Printf("(%d, %d) -> (%d, %d)\n", min_x, min_y, max_x, max_y)
			// Send only their visible part of the map
			Msg{ Type: PayTypPlot, Count: int32((1 + max_x - min_x) * (1 + max_y - min_y)) }.Write(enc)
			for x := min_x; x <= max_x; x += 1 {
				for y := min_y; y <= max_y; y += 1 {
					plot:=&GameWorld.Plots[x][y]

					if err := plot.Write(enc); err != nil {
						fmt.Println("Sending map failed!")
					}
				}
			}

			// player view and cam setup
			player.Cam.X, player.Cam.Y = WorldWidth / 2, WorldHeight / 2
			player.Cursor.X, player.Cursor.Y = player.Cam.X, player.Cam.Y
			player.Cam.View.Width = ViewWidth
			player.Cam.View.Height = ViewHeight
			player.MoveCursorTo(start.X, start.Y)

			// Send the initial empty player object
			if err := msg.Write(enc); err != nil { panic(err) }
			if err := player.Write(enc); err != nil { panic(err) }

			// Continuous message handling
			for {
				msg.Read(dec)

				switch {
				case msg.Type == PayTypJoin:
					player.Read(dec)
					players = append(players, player)
					fmt.Println(players[len(players)-1].Name + " has joined the game")

					// Spawn their village
					start.SpawnUnit(UnitVillage, &player)
					Msg{ Type: PayTypPlot, Count: 1 }.Write(enc)
					start.Write(enc)
				}

				msg.Type = -1
				// fmt.Println("Got something")
			}

			conn.Close()
		}()
	}
}
