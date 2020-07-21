package cursor

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/KlyuchnikovV/cui/types"
)

type Cursor struct {
	// Unused for now
	x, y           int
	cursorPosition *regexp.Regexp

	types.ConsoleStream
}

func New() (*Cursor, error) {
	r, err := regexp.Compile(cursorPositionFormat)
	if err != nil {
		return nil, err
	}
	c := &Cursor{
		cursorPosition: r,
		ConsoleStream:  types.NewConsoleStream(),
	}
	c.x, c.y, err = c.GetCursor()

	return c, err
}

func (c *Cursor) Move(where CursorMovement, nTimes int) {
	c.Print(fmt.Sprintf(where.getString(), nTimes))
}

func (c *Cursor) SetCursor(x, y int) {
	c.Print(fmt.Sprintf(setCursorPosition, x, y))
}

func (c *Cursor) GetCursor() (int, int, error) {
	c.Print(getCursorPosition)

	var bytes = make([]byte, 7)
	if _, err := os.Stdin.Read(bytes); err != nil {
		return 0, 0, err
	}

	return c.parsePosition(bytes)
}

func (c *Cursor) parsePosition(bytes []byte) (int, int, error) {
	matches := c.cursorPosition.FindStringSubmatch(string(bytes))
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("couldn't parse response: %s", string(bytes))
	}
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

// func (c *ConsoleUI) PrintAt(x, y int, s string) {
// 	c.Print(saveCurPos)
// 	c.SetCursor(x, y)
// 	c.Print(s)
// 	c.Print(restoreCurPos)
// }

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
