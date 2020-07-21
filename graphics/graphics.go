package graphics

import (
	"fmt"
	"strings"

	"github.com/KlyuchnikovV/cui/cursor"
	"github.com/KlyuchnikovV/cui/low_level/terminal"
	"github.com/KlyuchnikovV/cui/types"
)

type Graphics struct {
	types.ConsoleStream
	cursor.Cursor
	width  int
	height int
}

func New() (*Graphics, error) {
	c, err := cursor.New()
	if err != nil {
		return nil, err
	}

	w, h := terminal.GetTerminalSize()
	return &Graphics{
		ConsoleStream: types.NewConsoleStream(),
		Cursor:        *c,
		width:         w,
		height:        h,
	}, nil
}

func (g *Graphics) ClearScreen(mode ClearMode) {
	g.Print(fmt.Sprintf(clearScreen, mode.getModeInt()))
}

func (g *Graphics) ClearLine(mode clearLineMode) {
	g.Print(fmt.Sprintf(clearLine, mode.getLineModeInt()))
}

func (g *Graphics) SetGraphics(modes ...GraphicsMode) {
	var result = make([]string, len(modes))
	for i, mode := range modes {
		result[i] = mode.getGraphicsModeString()
	}
	g.Print(fmt.Sprintf(setGraphics, strings.Join(result, ";")))
}

func (g *Graphics) SetForegroundColor(color Color) {
	g.Print(fmt.Sprintf(setGraphics, setForegroundColor.addColorOffset(color)))
}

func (g *Graphics) SetBackgroundColor(color Color) {
	g.Print(fmt.Sprintf(setGraphics, setBackgroundColor.addColorOffset(color)))
}

func (g *Graphics) ResetForegroundColor() {
	g.Print(fmt.Sprintf(setGraphics, ResetForegroundColor.getGraphicsModeString()))
}

func (g *Graphics) ResetBackgroundColor() {
	g.Print(fmt.Sprintf(setGraphics, ResetBackgroundColor.getGraphicsModeString()))
}

func (g *Graphics) PrintAt(x, y int, s string, restorePosition bool) {
	if restorePosition {
		g.SavePosition()
		defer g.RestorePosition()
	}
	g.SetCursor(x, y)
	g.Print(s)
}

// func (g *Graphics) DrawRectangle(x, y, width, height int, symbol rune) error {
// 	// if x < 0 || x > c.width {
// 	// 	return fmt.Errorf("wrong x coordinate")
// 	// }

// 	var lines = make([]string, height)
// 	lines[0] = strings.Repeat(string(symbol), width)
// 	lines[len(lines)-1] = lines[0]
// 	for i := 1; i < len(lines)-1; i++ {
// 		lines[i] = fmt.Sprintf("%c%s%c", symbol, strings.Repeat(" ", width-2), symbol)
// 	}
// 	g.PrintAt(x, y, strings.Join(lines, "\n"))
// 	return nil
// }
