package main

import (
	"encoding/gob"
	"fmt"
)

const (
	PayTypText    int8 = 0
	PayTypJoin    int8 = 1
	PayTypPlayer  int8 = 2
	PayTypPlot    int8 = 3
	PayTypChat    int8 = 4
)

type Msg struct {
	Version int16
	Type int8
	Count int32
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


func (p *Player) Read(dec *gob.Decoder) error {
	return dec.Decode(p)
}

func (p *Player) Write(enc *gob.Encoder) error {
	return enc.Encode(p)
}
