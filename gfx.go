package main

import (
	"github.com/nsf/termbox-go"
	"fmt"
)

type MsgContainer struct {
	Str string
	Y int
}

type Menu struct {
	title string
	options []string
	callback func(int)
}

var gfxMsgQueue [10]MsgContainer
var gfxMenus [10]Menu
var gfxCurrentMenu Menu
var MenuQueue = gfxMenus[0:0]
var MsgQueue = gfxMsgQueue[0:0]
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


func GfxMsgExplicit(m MsgContainer) {
	MsgQueue = append(MsgQueue, m)
}


func GfxMsg(m string) {
	MsgQueue = append(MsgQueue, MsgContainer{ Str: m, Y: -1 })
}


func gfxStringAt(x, y int, m string) {
	m = " " + m + " "
	for i := 0; i < len(m); i+= 1 {
		termbox.SetCell(x + i, y, rune(m[i]), termbox.ColorWhite, termbox.ColorBlack)
	}
}


func gfxLineFromString(s string) string {
	for i := 0; i < len(s); i += 1 {
		if s[i] == '\n' { return fmt.Sprintf("%v", s[:i]) }
	}

	return s[:]
}


func gfxStringCenteredAt(y int, m string) {
	width, _ := termbox.Size()
	off := 0

	for i := 0; i < 2; i += 1 {
		line := gfxLineFromString(m[off:])
		w := (width - len(line)) / 2

		// fmt.Println(line)

		y += 1
		off += len(line) + 1

		gfxStringAt(w, y, line)
		if off >= len(m) {
			break
		}
	}
}


func gfxShowMsgs() {
	w, h := termbox.Size()
	m := MsgQueue[0]

	w = (w - len(m.Str)) / 2
	h /= 2

	y := h

	if m.Y > -1 {
		y = m.Y
	}

	gfxStringCenteredAt(y, m.Str)
	// gfxStringAt(w, h, m)

	termbox.Flush()

	evt := termbox.PollEvent()
	if evt.Type == termbox.EventInterrupt {
		return
	}

	MsgQueue = MsgQueue[1:len(MsgQueue)]
}


func gfxShowMenus() {
	menu := MenuQueue[0]
	_, th := termbox.Size()
	th >>= 1

	gfxStringCenteredAt(th + 1, menu.title)
	for i := 0; i < len(menu.options); i += 1 {
		gfxStringCenteredAt(th + i + 2, fmt.Sprintf("%d %v", i + 1, menu.options[i]))
	}

	termbox.Flush()

	evt := termbox.PollEvent()
	if evt.Type == termbox.EventInterrupt {
		return
	}

	i := int(evt.Ch - '0') - 1
	if i >= 0 && i < len(menu.options) {
		menu.callback(i)
	}

	MenuQueue = MenuQueue[1:len(MenuQueue)]
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


func GfxMenu(title string, options []string, on_selection func(int)) {
	MenuQueue = append(MenuQueue, Menu {
		title: title,
		options: options,
		callback: on_selection,
	})
}


func (w *World) GfxDraw(player *Player) {
	unitNone := UnitIndex("vacant")
	tw, _ := termbox.Size()
	thw := tw >> 1
	width, height := ViewWidth, ViewHeight
	hw, hh := width >> 1, height >> 1
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			_x := player.Cam.X - hw + x
			_y := player.Cam.Y - hh + y

			fg, bg := termbox.ColorWhite, termbox.ColorBlack
			symbol := rune('?')
			var plot *Plot

			if _x >= 0 && _y >= 0 && _x < GameWorld.Width && _y < GameWorld.Height {
				plot = &w.Plots[_x][_y]
			}

			if plot != nil && plot.Explored & player.ID > 0 {
				if plot.Unit.Type != unitNone {
					// Man made
					player := PlayerFromID(plot.Unit.OwnerID)
					bg = player.Colors.Bg
					fg = player.Colors.Fg | plot.Unit.Attr
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
		}
	}

	gfxStringCenteredAt(0, fmt.Sprintf("%v(%d) $: %0.2f | %0.2f $/sec", player.Name, player.Score, player.Wealth, player.Income))

	selPlot := player.SelectedPlot(w)
	if selPlot.Explored & player.ID > 0 {
		gfxStringCenteredAt(1, selPlot.Description(UnitDescriptionShort))
		// gfxStringCenteredAt(3, fmt.Sprintf("unit:%d owner:%d", selPlot.Unit.Type, selPlot.Unit.OwnerID))

		if selPlot.Unit.Type != unitNone {
			resources := selPlot.Unit.Resources
			gfxStringCenteredAt(2, fmt.Sprintf("$: %0.2f | %0.2f $/sec", resources.Current, resources.Rate * selPlot.Productivity))
		}
	}
}


func GfxDrawBegin() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
}


func GfxDrawFinish(poll bool) termbox.Event {
	var evt termbox.Event

	if len(MsgQueue) > 0 {
		gfxShowMsgs()
	} else if len(MenuQueue) > 0 {
		gfxShowMenus()
	} else{
		termbox.Flush()

		if poll {
			evt = termbox.PollEvent()
		}
	}

	return evt
}
