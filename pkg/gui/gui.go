package gui

import (
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	FG   = termbox.ColorDefault
	BG   = termbox.ColorDefault
	SepX = 11
	SepY = 3
)

var (
	MenuOptions = []string{
		"New Message",
		"New Contact",
		"Placeholder",
	}
)

type GUI struct {
	x    int
	line int

	menu *Menu
}

type Menu struct {
	options  []string
	selected int
}

func New() (*GUI, error) {
	if err := termbox.Init(); err != nil {
		return nil, fmt.Errorf("failed to init GUI: %v", err)
	}

	g := &GUI{
		x:    0,
		line: 0,
	}

	g.menu = &Menu{
		options:  MenuOptions,
		selected: 0,
	}

	g.init()

	return g, nil
}

func (g *GUI) init() {
	termbox.Clear(FG, BG)
	termbox.Flush()
	g.catchKeyboard()
}

func (g *GUI) Close() {
	termbox.Close()
}

func (g *GUI) DrawMenu() {
	g.line = 0
	termbox.Clear(FG, BG)
	termbox.Flush()
	g.drawText("go-whisper", 0, 1, FG, BG)
	w, h := termbox.Size()
	g.fill(SepX, 0, 1, h, termbox.Cell{Ch: '|'})
	g.fill(0, SepY, w, 1, termbox.Cell{Ch: '-'})

	x := SepX + 1
	for i, o := range g.menu.options {
		color := termbox.ColorDefault
		if i == g.menu.selected {
			color = termbox.ColorRed
		}

		x, _ = g.drawText(o, x+2, 1, FG, color)
	}

	//for i := uint16(0); i < 16; i++ {
	//	x, _ = g.drawText(fmt.Sprintf("%v", i), x+2, 1, FG, termbox.Attribute(i))
	//}

	//x, _ = g.drawText(fmt.Sprintf("%v %v", termbox.ColorBlue, termbox.ColorGreen), x+2, 1, FG, termbox.ColorBlue)

	termbox.Flush()
}

func (g *GUI) Print(msg string) {
	for _, c := range msg {
		if c == '\n' {
			g.line++
			continue
		}
		termbox.SetCell(g.x, g.line, c, FG, BG)
		g.x += runewidth.RuneWidth(c)
	}

	g.x = 0
	g.line++
	termbox.Flush()
}

func (g *GUI) drawText(msg string, x, y int, fg, bg termbox.Attribute) (int, int) {
	for _, c := range msg {
		if c == '\n' {
			y++
			continue
		}

		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
	termbox.Flush()

	return x, y
}

func (g *GUI) Infof(msg string) {
	g.Print(msg)
}

func (g *GUI) fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func (g *GUI) catchKeyboard() {
	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:

				switch ev.Key {
				case termbox.KeyCtrlC:
					termbox.Close()
					fmt.Printf("closing...\n")
					os.Exit(0)
					break

				case termbox.KeyArrowLeft:
					g.menu.selected = g.menu.selected - 1
					if g.menu.selected < 0 {
						g.menu.selected = len(g.menu.options) - 1
					}

					g.DrawMenu()
					break

				case termbox.KeyArrowRight:
					g.menu.selected = (g.menu.selected + 1) % len(g.menu.options)

					g.DrawMenu()
					break
				}

				break
			}
		}
	}()
}
