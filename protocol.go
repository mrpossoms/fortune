package main

import (
	"encoding/gob"
	"net"
)

const (
	PayTypText int8 = 0
	PayTypJoin      = 1
	PayTypPlayers   = 2
	PayTypPlot      = 3
	PayTypChat      = 4
)

type Msg struct {
	Version int16
	Type int8
	Count int32
}

func (m Msg) Write(con net.Conn) error {
	enc := gob.NewEncoder(con)
	return enc.Encode(m)
}


func (m *Msg) Read(con net.Conn) error {
	dec := gob.NewDecoder(con)
	return dec.Decode(m)
}


func (p *Plot) Read(con net.Conn) error {
	dec := gob.NewDecoder(con)
	return dec.Decode(p)
}


func (p Plot) Write(con net.Conn) error {
	enc := gob.NewEncoder(con)
	return enc.Encode(p)
}


func (w *World) Write(con net.Conn) error {
	for x := 0; x < WorldWidth; x += 1 {
		for y := 0; y < WorldHeight; y += 1 {
			err := w.Plots[x][y].Write(con)

			if err != nil { return err }
		}
	}

	return nil
}


func (p *Player) Read(con net.Conn) error {
	dec := gob.NewDecoder(con)
	return dec.Decode(p)
}

func (p Player) Write(con net.Conn) error {
	enc := gob.NewEncoder(con)
	return enc.Encode(p)
}
