package main

import (
	"os"
	"time"
	"os/user"
)

func Intro() {
	GfxInit()

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	GfxMsg("F")
	GfxMsg("FTE")
	GfxMsg("FRTNE")
	GfxMsg("FORTUNE")
	GfxMsg("F O R T U N E")

	for ;len(MsgQueue) > 0; {
		GfxDrawBegin()
		GfxDrawFinish(false)
		time.Sleep(100 * time.Millisecond)
	}
	GfxDrawFinish(true)

	if _, err := os.Stat(usr.HomeDir + "/.fortune"); os.IsNotExist(err) {
		GfxMsg("Welcome to Fortune!")
		GfxMsg("This is the first time you've been here!")
		GfxMsg("How about a little introduction?")
		GfxMsg("Fortune is a multiplayer empire building simulator")

		for ;len(MsgQueue) > 0; {
			GfxDrawBegin()
			GfxDrawFinish(true)
		}

		file, err := os.Create(usr.HomeDir + "/.fortune")
		if err != nil { panic(err) }
		file.Close()
	}


}
