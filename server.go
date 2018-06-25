package main

import (
	"os"
	"strconv"
	"time"
	"net"
	"fmt"
	"encoding/gob"
)

func leader(players []Player) (l *Player) {
	for pi := 0; pi < len(players); pi += 1 {
		if l == nil || players[pi].Score > l.Score {
			l = &players[pi]
		}
	}

	return
}

func GameServer(ln net.Listener) {
	var playerPool [64]Player

	players := playerPool[0:0]

	size, err := strconv.Atoi(os.Args[2])

	if err != nil {
		panic("Please enter a mapsize. Ex. 'fortune host 64'")
	}

	GameWorld = World {
		Width: size,
		Height: size,
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
	updateTicker := time.NewTicker(time.Second * 1)
	go func() {
		for {
			<-updateTicker.C
			GameWorld.Tick(gameTime)

			var plots [MaxWorldWidth * MaxWorldHeight]*Plot

			for conni := 0; conni < len(players); conni += 1 {
				pconn := PlayerConns[conni]

				changed := GameWorld.ChangedPlots(plots[0:0], gameTime, pconn.ID)
				pconn.Lock.Lock()
				Msg{ Type: PayTypPlot, Count: int32(len(changed)) }.Write(pconn.Enc)

				for plotIdx := 0; plotIdx < len(changed); plotIdx += 1 {
					changed[plotIdx].Write(pconn.Enc)
				}

				// send all player scores to each player
				//fmt.Printf("To: %d %v\n", pconn.ID, PlayerFromID(pconn.ID).Name)
				Msg{ Type: PayTypPlayer, Count: int32(len(players)) }.Write(pconn.Enc)
				for pi := 0; pi < len(players); pi += 1 {
					player := &players[pi]
					//fmt.Printf("\t%d %v\n", player.ID, player.Name)
					player.Wealth, player.Income = GameWorld.PlayerResources(player.ID)
					player.Score = int32(GameWorld.PlayerScore(player.ID))
					player.Write(pconn.Enc)
				}

				pconn.Lock.Unlock()
			}

			if GameWorld.IsGameOver() {
				winner := leader(players)
				for conni := 0; conni < len(players); conni += 1 {
					pconn := PlayerConns[conni]
					Msg{ Type: PayTypText, Count: 1 }.Write(pconn.Enc)
					TextPayload{ Msg: fmt.Sprintf("Game over, %s wins", winner.Name) }.Write(pconn.Enc)
				}

				fmt.Println("Game over")
				os.Exit(0)
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

			// Send game info to player
			Msg{ Type: PayTypInfo, Count:1 }.Write(pconn.Enc)
			GameInfo{ GameWorld.Width, GameWorld.Height, GameWorld.TotalCaptureSpace() }.Write(pconn.Enc)

			// Find a place for them to start
			start := GameWorld.FindLivablePlot()

			// player view and cam setup
			player.Cam.X, player.Cam.Y = GameWorld.Width / 2, GameWorld.Height / 2
			player.Cursor.X, player.Cursor.Y = player.Cam.X, player.Cam.Y
			player.Cam.View.Width = ViewWidth
			player.Cam.View.Height = ViewHeight
			player.MoveCursorTo(start.X, start.Y)

			// Send the initial empty player object
			if err := msg.Write(pconn.Enc); err != nil { panic(err) }
			if err := player.Write(pconn.Enc); err != nil { panic(err) }
			fmt.Println("Sent Player")

			// Continuous message handling
			for {
				pconn.Lock.Lock()
				msg.Read(pconn.Dec)
				switch {
				case msg.Type == PayTypJoin:
					player.Read(pconn.Dec)
					players = append(players, player)
					fmt.Printf("%v (%d) has joined the game\n", players[len(players)-1].Name, players[len(players)-1].ID)
					fmt.Printf("With color %d\n", players[len(players)-1].Colors.Bg)

					villageUi := UnitIndex("village")
					region := GameWorld.Reveal(start.X, start.Y, Units[villageUi].Sight, player.ID)

					// Spawn their village
					start.SpawnUnit(villageUi, &player)
					// Msg{ Type: PayTypPlot, Count: 1 }.Write(pconn.Enc)
					// start.Write(pconn.Enc)

					// Send only their visible part of the map
					Msg{ Type: PayTypPlot, Count: int32(region.Area()) }.Write(pconn.Enc)
					GameWorld.WriteRegion(pconn.Enc, region)


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
						// Deduct resources from neighboring plots
						unit := Units[updatedPlot.Unit.Type]
						cost := unit.Resources.Cost
						if GameWorld.Plots[x][y].SpendResources(&GameWorld, cost, player.ID) {
							GameWorld.Plots[x][y].SpawnUnit(unit.Type, &player)
							_ = GameWorld.Reveal(x, y, unit.Sight, player.ID)
						} else {
							Msg{ Type: PayTypText, Count: 1 }.Write(pconn.Enc)
							TextPayload{ Msg: fmt.Sprintf("Not enough resources to build a %s", unit.Name) }.Write(pconn.Enc)
							fmt.Println("Not enough resources")
						}
						// fmt.Printf("(%d, %d) -> (%d, %d)\n", min_x, min_y, max_x, max_y)

						// Send newly explored part of the map
						// Msg{ Type: PayTypPlot, Count: int32(region.Area()) }.Write(pconn.Enc)
						// GameWorld.WriteRegion(pconn.Enc, region)
					}

					break
				}

				msg.Type = -1
				pconn.Lock.Unlock()
				// fmt.Println("Got something")
			}

			conn.Close()
		}()
	}
}
