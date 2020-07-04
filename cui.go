package cui

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/KlyuchnikovV/chan_utils"
)

type ConsoleUI struct {
	ctx        context.Context
	cancel     context.CancelFunc
	outStream  io.Writer
	renderChan chan chan_utils.Message
}

func New(ctx context.Context, debug bool) (*ConsoleUI, error) {
	out, err := os.OpenFile("/dev/tty", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Printf("INFO: %s\n", err.Error())
		if debug {
			out = os.Stdout
		} else {
			return nil, err
		}
	}

	return &ConsoleUI{
		ctx:        ctx,
		outStream:  out,
		renderChan: make(chan chan_utils.Message, 10),
	}, nil
}

func (c *ConsoleUI) StartRendering(async bool) {
	if c.cancel != nil {
		return
	}

	render, cancel := chan_utils.NewListener(c.ctx, c.renderChan, c.onRenderRequest, func(err error) {
		log.Print(err)
		panic(err)
	})

	c.cancel = cancel

	if async {
		go render()
	} else {
		render()
	}
}

func (c *ConsoleUI) Stop() {
	if c.cancel == nil {
		return
	}

	c.cancel()
	c.cancel = nil
}

func (c *ConsoleUI) onRenderRequest(data chan_utils.Message) {
	msg, ok := data.(renderRequest)
	if !ok {
		panic("!ok")
	}
	log.Printf("INFO: got to render " + string(msg))
	_, err := c.outStream.Write(msg)
	if err != nil {
		panic(err)
	}
}

func (c *ConsoleUI) Write(p []byte) {
	c.outStream.Write(p)
}

func (c *ConsoleUI) Print(s string) {
	c.Write([]byte(s))
}

func (c *ConsoleUI) GetRenderChan() chan chan_utils.Message {
	return c.renderChan
}

func (c *ConsoleUI) Move(where CursorMovement, nTimes int) {
	c.Print(fmt.Sprintf(where.getString(), nTimes))
}

func (c *ConsoleUI) SetCursor(x, y int) {
	c.Print(fmt.Sprintf(setCursorPosition, x, y))
}

func (c *ConsoleUI) ClearScreen(mode ClearMode) {
	c.Print(fmt.Sprintf(clearScreen, mode.getModeInt()))
}

func (c *ConsoleUI) ClearLine(mode clearLineMode) {
	c.Print(fmt.Sprintf(clearLine, mode.getLineModeInt()))
}

func (c *ConsoleUI) PrintAt(x, y int, s string) {
	c.Print(saveCurPos)
	c.SetCursor(x, y)
	c.Print(s)
	c.Print(restoreCurPos)
}

func (c *ConsoleUI) ShowCursor(show bool) {
	var mode string
	if show {
		mode = showCursor
	} else {
		mode = hideCursor
	}
	c.Print(mode)
}

func (c *ConsoleUI) ScrollUp(lines int) {
	c.Print(fmt.Sprintf(scrollUp, lines))
}

func (c *ConsoleUI) ScrollDown(lines int) {
	c.Print(fmt.Sprintf(scrollDown, lines))
}

func (c *ConsoleUI) SetGraphics(modes ...GraphicsMode) {
	var result = make([]string, len(modes))
	for i, mode := range modes {
		result[i] = mode.getGraphicsModeString()
	}
	c.Print(fmt.Sprintf(setGraphics, strings.Join(result, ";")))
}

func (c *ConsoleUI) SetForegroundColor(color Color) {
	c.Print(fmt.Sprintf(setGraphics, setForegroundColor.addColorOffset(color)))
}

func (c *ConsoleUI) SetBackgroundColor(color Color) {
	c.Print(fmt.Sprintf(setGraphics, setBackgroundColor.addColorOffset(color)))
}

func (c *ConsoleUI) ResetForegroundColor() {
	c.Print(fmt.Sprintf(setGraphics, ResetForegroundColor.getGraphicsModeString()))
}

func (c *ConsoleUI) ResetBackgroundColor() {
	c.Print(fmt.Sprintf(setGraphics, ResetBackgroundColor.getGraphicsModeString()))
}
