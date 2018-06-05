package gui

import (
	"github.com/nsf/termbox-go"
)

type Contact struct {
	text []byte
	gui  *GUI

	stream chan rune
	stopCh chan struct{}
}

func newContact(gui *GUI, stream chan rune, stopCh chan struct{}) *Contact {
	contact := &Contact{
		stream: stream,
		stopCh: stopCh,
		gui:    gui,
	}

	contact.listenToSteam()

	return contact
}

func (c *Contact) printNewContact() {
	w, h := termbox.Size()
	headStr := "Please enter new user UID:"
	c.gui.drawText(headStr, (w-stringLength(headStr))/2, h/2-3, FG, BG)

	srtX, srtY := w/2-15, h/2-1
	endX, endY := w/2+15, h/2+1
	c.gui.fill(srtX, srtY, 30, 1, termbox.Cell{Ch: '─'})
	c.gui.fill(srtX, endY, 30, 1, termbox.Cell{Ch: '─'})

	termbox.SetCell(srtX-1, srtY, '┌', FG, BG)
	termbox.SetCell(srtX-1, endY, '└', FG, BG)
	termbox.SetCell(endX, srtY, '┐', FG, BG)
	termbox.SetCell(endX, endY, '┘', FG, BG)
	termbox.SetCell(srtX-1, srtY+1, '│', FG, BG)
	termbox.SetCell(endX, endY-1, '│', FG, BG)

	termbox.SetCursor(srtX+1, srtY+1)
	c.gui.drawText("Press enter to confirm choice.", srtX, h/2+2, FG, BG)

	termbox.Flush()
}

func (c *Contact) listenToSteam() {
	listen := func() bool {
		select {
		case ch := <-c.stream:
			c.gui.Print(string(ch))
			// insert text
			break

		case <-c.stopCh:
			return false

		}

		return true
	}

	go func() {
		for {
			if !listen() {
				break
			}
		}
	}()
}
