package main

import (
	"fmt"
	"net"
	// "time"
	"os"
)

var GameWorld World
var Players [64]Player
var PlayerConns [64]PlayerConnection

func main() {
	fmt.Println(os.Getpid())

	if os.Args[1] == "host" {
		ln, err := net.Listen("tcp", ":31337")

		if err != nil {
			panic(err.Error())
		}

		GameServer(ln)
	} else {
		// Intro()

		GameWorld = World {
			Smoothness: 2,
		}

		//Generate the landscape
		// for i := 0; i < 3; i += 1 {
		// 	GameWorld.Init(int64(time.Now().Second()))
		// }

		// Can't host, act as a client
		GameClient()
	}
}
