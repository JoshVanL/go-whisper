package gui

import (
	//"errors"
	//"strconv"

	"github.com/nsf/termbox-go"
)

type NewMsg struct {
	gui  *GUI
	uids []string
	text string

	cursorX, cursorY int
	startX, startY   int
	selected         int

	stream chan rune
	key    chan termbox.Key
	stopCh chan struct{}
}

func newMessage(gui *GUI, stream chan rune, key chan termbox.Key, stopCh chan struct{}) *NewMsg {
	newMsg := &NewMsg{
		stream:   stream,
		key:      key,
		stopCh:   stopCh,
		gui:      gui,
		uids:     gui.client.Uids(),
		selected: 0,
	}

	newMsg.listenToSteam()

	return newMsg
}

func (c *NewMsg) printNewMessage() {
	termbox.Flush()
	termbox.Clear(FG, BG)
	c.gui.DrawMenu()

	w, h := termbox.Size()
	headStr := "Select uid to send message to:"
	c.gui.drawText(headStr, (w-stringLength(headStr))/4, h/6-3, FG, BG)

	yy := h/6 - 1
	for s, uid := range c.uids {
		fg := termbox.ColorCyan
		bg := termbox.ColorBlack
		if s == c.selected {
			fg = termbox.ColorDefault
			bg = termbox.ColorRed
		}

		c.gui.drawText(uid, (w-stringLength(headStr))/4, yy, fg, bg)
		yy++
	}

}

func (n *NewMsg) listenToSteam() {
	listen := func() bool {
		select {
		case ch := <-n.stream:

			//n.gui.drawText(n.clearBoxString(30), n.startX, n.cursorY, FG, BG)

			if len(n.text) < 27 {
				n.text += string(ch)
				n.cursorX++
			}

			n.gui.drawText(n.text, n.startX, n.cursorY, FG, BG)
			termbox.SetCursor(n.cursorX, n.cursorY)
			n.gui.drawText(n.text, n.startX, n.cursorY, FG, BG)

			break

		case key := <-n.key:
			if key == termbox.KeyArrowUp {
				n.selected--
				if n.selected < 0 {
					n.selected = len(n.uids) - 1
				}

			} else if key == termbox.KeyArrowDown {
				n.selected = (n.selected + 1) % len(n.uids)

			} else if key == termbox.KeyTab {

				//n.gui.drawText(n.clearBoxString(100), n.startX-1, n.startY+4, FG, BG)

				//res, err := n.enterUid()
				//if err != nil {
				//	n.gui.drawText(err.Error(), n.startX-1, n.cursorY+4, FG, termbox.ColorRed)
				//	break
				//}

				//n.gui.drawText(res, n.startX-1, n.cursorY+4, FG, termbox.ColorCyan)
			}

			n.printNewMessage()

			break

		case <-n.stopCh:
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
