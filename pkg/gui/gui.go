package gui

import (
	"fmt"
	"os"
	"sync"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"

	"github.com/joshvanl/go-whisper/pkg/interfaces"
)

const (
	FG   = termbox.ColorDefault
	BG   = termbox.ColorDefault
	SepX = 13
	SepY = 3
)

var (
	MenuOptions = []string{
		"New Message",
		"New Contact",
		"Chats",
	}
)

type GUI struct {
	uid      uint64
	x        int
	line     int
	initMenu bool

	stream chan rune
	keys   chan termbox.Key

	stopPage  chan struct{}
	enterMode bool
	mu        *sync.Mutex

	menu    *Menu
	contact *Contact
	client  interfaces.Client
}

type Menu struct {
	options  []string
	selected int
	page     int
}

func New() (*GUI, error) {
	if err := termbox.Init(); err != nil {
		return nil, fmt.Errorf("failed to init GUI: %v", err)
	}

	g := &GUI{
		uid:      0,
		x:        0,
		line:     0,
		stream:   make(chan rune),
		keys:     make(chan termbox.Key),
		mu:       new(sync.Mutex),
		stopPage: make(chan struct{}),
	}

	g.menu = &Menu{
		options:  MenuOptions,
		selected: 0,
		page:     0,
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
	//if !g.initMenu {
	//termbox.Clear(FG, BG)
	g.initMenu = true
	//}

	g.line = 0
	//termbox.Flush()
	w, h := termbox.Size()
	g.fill(0, 0, w, h, termbox.Cell{Ch: ' '})
	g.drawText("go-whisper", 1, 1, termbox.ColorCyan, termbox.ColorBlack)
	//w, h := termbox.Size()
	g.fill(SepX, 0, 1, h, termbox.Cell{Ch: '|'})
	g.fill(0, SepY, w, 1, termbox.Cell{Ch: '-'})

	yy := 5
	for _, uid := range g.client.Uids() {
		g.drawText(fmt.Sprintf("%s", uid), 1, yy, termbox.ColorCyan, termbox.ColorBlack)
		yy++
	}

	uid := fmt.Sprintf("%v", g.uid)
	for len(uid) != 11 {
		uid = fmt.Sprintf("0%s", uid)
	}

	pageStr := fmt.Sprintf("%s uid[%s]", g.menu.options[g.menu.page], uid)
	g.drawText(pageStr, w-stringLength(pageStr)-1, 2, FG, termbox.ColorMagenta)

	x := SepX + 1
	for i, o := range g.menu.options {
		color := termbox.ColorDefault
		if i == g.menu.selected {
			color = termbox.ColorRed
		}

		x, _ = g.drawText(o, x+2, 1, FG, color)
	}

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
	keybord := func() {
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

			case termbox.KeyEnter:
				g.mu.Lock()
				termbox.SetCursor(termbox.Size())
				close(g.stopPage)
				close(g.stream)
				close(g.keys)
				g.stream = make(chan rune)
				g.keys = make(chan termbox.Key)
				g.stopPage = make(chan struct{})
				g.mu.Unlock()

				switch g.menu.selected {
				case 0:
					g.enterMode = false
					g.initMenu = false
					g.menu.page = 0
					g.DrawMenu()
					break

				case 1:
					g.enterMode = true
					g.menu.page = 1
					g.DrawMenu()
					g.contact = newContact(g, g.stream, g.keys, g.stopPage)
					g.contact.printNewContact()
					break

				case 2:
					g.enterMode = false
					g.initMenu = false
					g.menu.page = 2
					g.DrawMenu()
					break

				}

				break

			default:
				if g.enterMode {
					g.mu.Lock()
					if ev.Ch != 0 {
						go func() {
							g.stream <- ev.Ch
						}()
					} else if ev.Key != 0 {
						go func() {
							g.keys <- ev.Key
						}()
					}
					g.mu.Unlock()
				}

				break
			}

			break
		}
	}

	go func() {
		for {
			keybord()
		}
	}()
}

func stringLength(msg string) (x int) {
	for _, c := range msg {
		x += runewidth.RuneWidth(c)
	}
	return x
}

func (g *GUI) SetUid(uid uint64) {
	g.uid = uid
}

func (g *GUI) SetClient(client interfaces.Client) {
	g.client = client
}
