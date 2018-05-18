package main

import (
	"github.com/nsf/termbox-go"
)


var _gfxMsgQueue [10]string
var msgQueue = _gfxMsgQueue[0:0]

func GfxInit() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}
}


func GfxUninit() {
	termbox.Close()
}


func GfxMsg(m string) {
	msgQueue = append(msgQueue, m)
}

func gfxStringAt(x, y int, m string) {
	for i := 0; i < len(m); i+= 1 {
		termbox.SetCell(x + i, y, rune(m[i]), termbox.ColorWhite, termbox.ColorBlack)
	}
}

func gfxShowMsgs() {
	w, h := termbox.Size()
	m := msgQueue[len(msgQueue) - 1]

	w = (w - len(m)) / 2
	h /= 2

	gfxStringAt(w, h, m)

	msgQueue = msgQueue[:len(msgQueue) - 1]
}

func (w World) GfxDraw() {
	width, height := 80, 32
	hw, hh := width >> 1, height >> 1
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			_x := w.Cam.X - hw + x
			_y := w.Cam.Y - hh + y
			plot := &w.Plots[_x][_y]

			fg, bg := termbox.ColorWhite, termbox.ColorBlack
			symbol := rune(' ')

			if plot.Explored > 0 {
				switch {
				case plot.Elevation < 0:
					bg = termbox.ColorBlue
					fg = bg + 1
					symbol = rune('~')
					break
				case plot.Elevation < 0.125:
					bg = termbox.ColorYellow
					symbol = rune(' ')
					break
				case plot.Elevation < 0.5:
					fg = termbox.ColorBlack | termbox.AttrUnderline
					bg = termbox.ColorGreen
					symbol = rune('^')
					break
				case plot.Elevation > 1:
					termbox.SetOutputMode(termbox.OutputGrayscale)
					fg = termbox.ColorBlack
					bg = termbox.ColorWhite
					symbol = rune('^')
					break
				}
			}

			if _x == w.Cursor.X && _y == w.Cursor.Y {
				fg = termbox.ColorWhite
				bg = termbox.ColorBlack
				symbol = rune('X')
			}

			termbox.SetCell(x, y, symbol, fg, bg)
			termbox.SetOutputMode(termbox.Output256)
		}
	}
}

func GfxDrawBegin() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func GfxDrawFinish() termbox.Event {
	if len(msgQueue) > 0 {
		gfxShowMsgs()
	}

	termbox.Flush()
	evt := termbox.PollEvent()

	return evt
}
