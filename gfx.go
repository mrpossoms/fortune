package main

import (
	"github.com/nsf/termbox-go"
)


var _gfxMsgQueue [10]string
var MsgQueue = _gfxMsgQueue[0:0]
var initialized bool

func GfxInit() {
	if !initialized {
		err := termbox.Init()

		if err != nil {
			panic(err)
		}

		termbox.SetOutputMode(termbox.Output256)
		initialized = true
	}
}


func GfxUninit() {
	termbox.Close()
}


func GfxMsg(m string) {
	MsgQueue = append(MsgQueue, m)
}

func gfxStringAt(x, y int, m string) {
	m = " " + m + " "
	for i := 0; i < len(m); i+= 1 {
		termbox.SetCell(x + i, y, rune(m[i]), termbox.ColorWhite, termbox.ColorBlack)
	}
}

func gfxStringCenteredAt(y int, m string) {
	w, _ := termbox.Size()
	w = (w - len(m)) / 2

	gfxStringAt(w, y, m)
}

func gfxShowMsgs() {
	w, h := termbox.Size()
	m := MsgQueue[0]

	w = (w - len(m)) / 2
	h /= 2

	gfxStringAt(w, h, m)

	MsgQueue = MsgQueue[1:len(MsgQueue)]
}


func GfxPrompt(prompt string) string {
	input := ""
	_, h := termbox.Size()
	h /= 2

	gfxStringCenteredAt(h - 1, prompt + ", then press ENTER")
	for {
		termbox.Flush()
		evt := termbox.PollEvent()

		switch {
		case evt.Key == termbox.KeyBackspace:
		case evt.Key == termbox.KeyBackspace2:
			if len(input) > 0 {
				input = input[:len(input) - 1]
			}
			break
		case evt.Key == termbox.KeySpace:
			input += " "
			break
		case evt.Key == termbox.KeyEnter:
			return input
		case evt.Ch != 0:
			input += string(evt.Ch)
		}

		gfxStringCenteredAt(h, input)
	}

	gfxStringCenteredAt(h, input)
	termbox.Flush()

	return input
}


func (w *World) GfxDraw(player *Player) {
	tw, _ := termbox.Size()
	thw := tw >> 1
	width, height := ViewWidth, ViewHeight
	hw, hh := width >> 1, height >> 1
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			_x := player.Cam.X - hw + x
			_y := player.Cam.Y - hh + y
			plot := &w.Plots[_x][_y]

			fg, bg := termbox.ColorWhite, termbox.ColorBlack
			symbol := rune('?')

			if plot.Explored & player.ID > 0 {
				if plot.Unit.Type != UnitNone {
					// Man made
					bg = termbox.ColorRed //plot.Unit.Owner.Colors.Bg
					fg = termbox.ColorWhite //plot.Unit.Owner.Colors.Fg | plot.Unit.Attr
					symbol = plot.Unit.Symbol
				} else {
					// Natural, unoccupied plot
					switch {
					case plot.Elevation < PlotSea:
						bg = termbox.ColorBlue
						fg = bg + 1
						symbol = rune('~')
						break
					case plot.Elevation < PlotBeach:
						bg = termbox.ColorYellow
						symbol = rune('¨')
						break
					case plot.Elevation < PlotPlains:
						fg = termbox.ColorBlack
						bg = termbox.ColorGreen
						symbol = rune('¸')
						break
					case plot.Elevation < PlotForest:
						fg = termbox.ColorBlack | termbox.AttrUnderline
						bg = termbox.ColorGreen
						symbol = rune('^')
						break
					case plot.Elevation < PlotMountain:
						// termbox.SetOutputMode(termbox.OutputGrayscale)
						fg = termbox.ColorBlack
						bg = termbox.ColorWhite
						symbol = rune('^')
						break
					}
				}

			}

			if _x == player.Cursor.X && _y == player.Cursor.Y {
				fg = termbox.ColorWhite
				bg = termbox.ColorBlack
				symbol = rune('X')
			}

			termbox.SetCell(x + thw - hw, y, symbol, fg, bg)
			termbox.SetOutputMode(termbox.Output256)
		}
	}

	if player.SelectedPlot(w).Explored & player.ID > 0 {
		gfxStringCenteredAt(1, player.SelectedPlot(w).Description(UnitDescriptionShort))
	}
}


func GfxDrawBegin() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
}


func GfxDrawFinish(poll bool) termbox.Event {
	if len(MsgQueue) > 0 {
		gfxShowMsgs()
	}

	termbox.Flush()

	var evt termbox.Event

	if poll {
		evt = termbox.PollEvent()
	}

	return evt
}
