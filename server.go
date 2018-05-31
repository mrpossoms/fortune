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

	// Game update
	gameTime := 0
	go func() {
		updateTicker := time.NewTicker(time.Second)

		for {
			<-updateTicker.C
			GameWorld.Tick(gameTime)

			var plots [WorldWidth * WorldHeight]*Plot

			for pi := 0; pi < len(players); pi += 1 {
				pconn := PlayerConns[pi]
				changed := GameWorld.ChangedPlots(plots[0:0], gameTime, pconn.ID)

				Msg{ Type: PayTypPlot, Count: int32(len(changed)) }.Write(pconn.Enc)

				for ci := 0; ci < len(changed); ci += 1 {
					changed[ci].Write(pconn.Enc)
				}
			}

			gameTime += 1
		}
	}()

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
					fmt.Printf("%v (%d) has joined the game\n", players[len(players)-1].Name, players[len(players)-1].ID)

					region := GameWorld.Reveal(start.X, start.Y, 3, player.ID)

					// Send only their visible part of the map
					Msg{ Type: PayTypPlot, Count: int32(region.Area()) }.Write(pconn.Enc)
					// for x := min_x; x <= max_x; x += 1 {
					// 	for y := min_y; y <= max_y; y += 1 {
					// 		plot:=&GameWorld.Plots[x][y]
					//
					// 		if err := plot.Write(pconn.Enc); err != nil {
					// 			fmt.Println("Sending map failed!")
					// 		}
					// 	}
					// }
					GameWorld.WriteRegion(pconn.Enc, region)

					// Spawn their village
					start.SpawnUnit(UnitVillage, &player)
					Msg{ Type: PayTypPlot, Count: 1 }.Write(pconn.Enc)
					start.Write(pconn.Enc)

					Msg{ Type: PayTypPlayer, Count: int32(len(players))}.Broadcast(players)
					for i := 0; i < len(players); i += 1 {
						players[i].Broadcast(players)
					}

				case msg.Type == PayTypPlot:
					updatedPlot := Plot{}
					updatedPlot.Read(pconn.Dec)
					x, y := updatedPlot.X, updatedPlot.Y
					plot := GameWorld.Plots[x][y]

					// Make sure that unit is allowed to be built there
					buildables := plot.PossibleBuilds(&GameWorld, pconn.ID)
					isBuildableUnit := false
					for i := 0; i < len(buildables); i += 1 {
						if updatedPlot.Unit.Type == buildables[i] {
							isBuildableUnit = true
							break
						}
					}

					fmt.Printf("Got plot (%d,%d)\n", x, y)
					if isBuildableUnit {
						GameWorld.Plots[x][y] = updatedPlot

						region := GameWorld.Reveal(x, y, 2, player.ID)
						// fmt.Printf("(%d, %d) -> (%d, %d)\n", min_x, min_y, max_x, max_y)

						// Send newly explored part of the map
						Msg{ Type: PayTypPlot, Count: int32(region.Area()) }.Write(pconn.Enc)
						// for x := min_x; x <= max_x; x += 1 {
						// 	for y := min_y; y <= max_y; y += 1 {
						// 		plot:=&GameWorld.Plots[x][y]
						//
						// 		if err := plot.Write(pconn.Enc); err != nil {
						// 			fmt.Println("Sending map failed!")
						// 		}
						// 	}
						// }
						GameWorld.WriteRegion(pconn.Enc, region)
					}

					break
				}

				msg.Type = -1
				// fmt.Println("Got something")
			}

			conn.Close()
		}()
	}
}
