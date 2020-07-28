package cursor

import (
	"fmt"

	"github.com/KlyuchnikovV/cui/types"
)

// TODO: how to automatically test GetCursor

type Cursor struct {
	types.ConsoleStream
}

func New() *Cursor {
	return &Cursor{
		ConsoleStream: types.NewConsoleStream(),
	}
}

func (c *Cursor) SetCursor(x, y int) {
	c.Print(fmt.Sprintf(setCursorPosition, x, y))
}

func (c *Cursor) GetCursor() (int, int, error) {
	c.Print(getCursorPosition)

	var x, y int
	_, err := fmt.Scanf(cursorPositionFormat, &x, &y)
	return x, y, err
}

func (c *Cursor) Move(where CursorMovement, nTimes int) {
	c.Print(fmt.Sprintf(where.getString(), nTimes))
}

func (c *Cursor) ShowCursor(show bool) {
	var mode string
	if show {
		mode = showCursor
	} else {
		mode = hideCursor
	}
	c.Print(mode)
}

func (c *Cursor) ScrollUp(lines int) {
	c.Print(fmt.Sprintf(scrollUp, lines))
}

func (c *Cursor) ScrollDown(lines int) {
	c.Print(fmt.Sprintf(scrollDown, lines))
}

func (c *Cursor) SavePosition() {
	c.Print(saveCursorPosition)
}

func (c *Cursor) RestorePosition() {
	c.Print(restoreCursorPosition)
}
