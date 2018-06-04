package gui

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const preferred_horizontal_threshold = 5
const tabstop_length = 8
const fg = termbox.ColorDefault
const bg = termbox.ColorDefault

type GUI struct {
	text []byte

	x    int
	line int
}

func New() (*GUI, error) {
	if err := termbox.Init(); err != nil {
		return nil, fmt.Errorf("failed to init GUI: %v", err)
	}

	g := &GUI{
		x:    0,
		line: 0,
	}

	g.init()

	return g, nil
}

func (g *GUI) init() {
	termbox.Clear(fg, bg)
	termbox.Flush()
	g.catch()
}

func (g *GUI) Close() {
	termbox.Close()
}

func (g *GUI) Print(msg string) {
	for _, c := range msg {
		termbox.SetCell(g.x, g.line, c, fg, bg)
		g.x += runewidth.RuneWidth(c)
	}
	termbox.SetCell(g.x, g.line, '\n', fg, bg)

	g.x = 0
	g.line++
	termbox.Flush()
}

func (g *GUI) Infof(msg string) {
	for _, c := range msg {
		termbox.SetCell(g.x, g.line, c, fg, bg)
		g.x += runewidth.RuneWidth(c)
	}
	termbox.SetCell(g.x, g.line, '\n', fg, bg)

	g.x = 0
	g.line++
	termbox.Flush()
}

func (g *GUI) catch() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlC {
					termbox.Close()
					fmt.Printf("closing...\n")
					os.Exit(0)
					break
				}
			}
		}
	}()
}
