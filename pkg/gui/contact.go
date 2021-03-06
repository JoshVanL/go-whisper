package gui

import (
	"errors"
	"strconv"

	"github.com/nsf/termbox-go"
)

type Contact struct {
	text string
	gui  *GUI

	cursorX, cursorY int
	startX, startY   int

	stream chan rune
	key    chan termbox.Key
	stopCh chan struct{}
}

func newContact(gui *GUI, stream chan rune, key chan termbox.Key, stopCh chan struct{}) *Contact {
	contact := &Contact{
		stream: stream,
		key:    key,
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

	c.cursorX = srtX + 1
	c.cursorY = srtY + 1
	c.startX = srtX + 1
	c.startY = srtY + 1

	termbox.SetCursor(c.cursorX, c.cursorY)
	c.gui.drawText("Press TAB to confirm choice.", srtX, h/2+2, FG, BG)

	termbox.Flush()
}

func (c *Contact) listenToSteam() {
	listen := func() bool {
		select {
		case ch := <-c.stream:

			c.gui.drawText(c.clearBoxString(30), c.startX, c.cursorY, FG, BG)

			if len(c.text) < 27 {
				c.text += string(ch)
				c.cursorX++
			}

			c.gui.drawText(c.text, c.startX, c.cursorY, FG, BG)
			termbox.SetCursor(c.cursorX, c.cursorY)
			c.gui.drawText(c.text, c.startX, c.cursorY, FG, BG)

			break

		case key := <-c.key:
			if key == termbox.KeyBackspace || key == termbox.KeyBackspace2 {
				if len(c.text) > 0 {
					c.text = c.text[:len(c.text)-1]
					c.cursorX--
					c.gui.drawText(c.clearBoxString(30), c.startX, c.cursorY, FG, BG)
					termbox.SetCursor(c.cursorX, c.cursorY)
					c.gui.drawText(c.text, c.startX, c.cursorY, FG, BG)
				}
			} else if key == termbox.KeyTab {

				c.gui.drawText(c.clearBoxString(100), c.startX-1, c.startY+4, FG, BG)

				res, err := c.enterUid()
				if err != nil {
					c.gui.drawText(err.Error(), c.startX-1, c.cursorY+4, FG, termbox.ColorRed)
					break
				}

				c.gui.drawText(res, c.startX-1, c.cursorY+4, FG, termbox.ColorCyan)
			}

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

func (c *Contact) clearBoxString(clear int) string {
	str := ""
	for i := 0; i < clear-1; i++ {
		str += " "
	}
	return str
}

func (c *Contact) enterUid() (string, error) {

	if _, err := strconv.Atoi(c.text); err != nil || len(c.text) != 11 {
		return "", errors.New("UIDs must be 11 digits.")
	}

	res, err := c.gui.client.QueryUID(c.text)
	if err != nil {
		return "", err
	}

	return res, nil
}
