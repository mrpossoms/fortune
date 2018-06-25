package main

import (
	"encoding/gob"
	"sync"
	"net"
	"fmt"
)

const (
	PayTypJoin    int8 = 1
	PayTypPlayer  int8 = 2
	PayTypPlot    int8 = 3
	PayTypChat    int8 = 4
	PayTypText    int8 = 5
	PayTypInfo    int8 = 6
)

type Msg struct {
	Version int16
	Type int8
	Count int32
}


type PlayerConnection struct {
	Conn net.Conn
	Index int
	ID int64
	Lock sync.Mutex

	Enc *gob.Encoder
	Dec *gob.Decoder
}


type TextPayload struct {
	Msg string
}


type GameInfo struct {
	Width, Height int
	CaptureSpace int
}


func (m GameInfo) Write(enc *gob.Encoder) error {
	return enc.Encode(m)
}

func (m *GameInfo) Read(dec *gob.Decoder) error {
	return dec.Decode(m)
}


func (m TextPayload) Write(enc *gob.Encoder) error {
	return enc.Encode(m)
}

func (m *TextPayload) Read(dec *gob.Decoder) error {
	return dec.Decode(m)
}


func (m Msg) Broadcast(players []Player) error {
	var err error

	for i := 0; i < len(players); i += 1 {
		idx := PlayerIndex(players[i].ID)
		err = m.Write(PlayerConns[idx].Enc)
	}

	return err
}


func (m Msg) Write(enc *gob.Encoder) error {
	return enc.Encode(m)
}


func (m *Msg) Read(dec *gob.Decoder) error {
	return dec.Decode(m)
}


func (p *Plot) Read(dec *gob.Decoder) error {
	return dec.Decode(p)
}


func (p *Plot) Write(enc *gob.Encoder) error {
	return enc.Encode(p)
}


func (w *World) Write(enc *gob.Encoder) error {
	for x := 0; x < len(w.Plots); x += 1 {
		for y := 0; y < len(w.Plots[x]); y += 1 {
			err := w.Plots[x][y].Write(enc)
			if err != nil {
				fmt.Println("World Write error")
				return err
			}
		}
	}

	return nil
}

func (w *World) WriteRegion(enc *gob.Encoder, region Region) error {
	for x := region.Min.X; x <= region.Max.X; x += 1 {
		for y := region.Min.Y; y <= region.Max.Y; y += 1 {
			plot:=&w.Plots[x][y]

			if err := plot.Write(enc); err != nil {
				fmt.Println("Sending map failed!")
			}
		}
	}

	return nil
}


func (p *Player) Read(dec *gob.Decoder) error {
	return dec.Decode(p)
}


func (p Player) Write(enc *gob.Encoder) error {
	return enc.Encode(p)
}


func (p *Player) Broadcast(players []Player) error {
	var err error

	for i := 0; i < len(players); i += 1 {
		idx := PlayerIndex(players[i].ID)
		err = p.Write(PlayerConns[idx].Enc)
	}

	return err
}
