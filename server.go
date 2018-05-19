package main

import (
	"time"
	"net"
	"fmt"
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

		fmt.Println("Connection incoming")
		go func() {
			msg := Msg{ Type: PayTypJoin, Count:1 }
			player := Player{ ID: 1 << uint(len(players))}

			// player view and cam setup
			player.Cam.X, player.Cam.Y = WorldWidth / 2, WorldHeight / 2
			player.Cursor.X, player.Cursor.Y = player.Cam.X, player.Cam.Y
			player.Cam.View.Width = ViewWidth
			player.Cam.View.Height = ViewHeight

			// Spawn their village
			// start:=GameWorld.FindLivablePlot()
			// start.SpawnUnit(UnitVillage, &player)
			// player.MoveCursorTo(start.X, start.Y)

			// Send the player object
			msg.Write(conn)
			player.Write(conn)

			// Continuous message handling
			for {
				msg.Read(conn)

				switch {
				case msg.Type == PayTypJoin:
					player.Read(conn)
					players = append(players, player)
					fmt.Println(players[len(players)-1].Name + " has joined the game")

					// Send the entire map
					Msg{ Type: PayTypPlot, Count:WorldWidth * WorldHeight }.Write(conn)
					GameWorld.Write(conn)
					fmt.Println("Sent map")
				}

				msg.Type = -1
				// fmt.Println("Got something")
			}

			conn.Close()
		}()
	}
}
