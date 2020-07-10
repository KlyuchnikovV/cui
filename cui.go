package cui

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/KlyuchnikovV/chan_utils"
	"github.com/KlyuchnikovV/cui/low_level/raw_mode"
)

type ConsoleUI struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	outStream           io.ReadWriter
	renderChan          chan chan_utils.Message
	cursorPositionRegex *regexp.Regexp
}

func New(ctx context.Context, enableRaw, debug bool) (*ConsoleUI, error) {
	out, err := os.OpenFile("/dev/tty", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Printf("INFO: %s\n", err.Error())
		if debug {
			out = os.Stdout
		} else {
			return nil, err
		}
	}

	if enableRaw {
		log.Print("entering raw")
		raw_mode.EnableRawMode()
	}

	r, err := regexp.Compile(cursorPositionFormat)
	if err != nil {
		return nil, err
	}

	return &ConsoleUI{
		ctx:                 ctx,
		outStream:           out,
		renderChan:          make(chan chan_utils.Message, 10),
		cursorPositionRegex: r,
	}, nil
}

func (c *ConsoleUI) DisableRawMode() {
	raw_mode.DisableRawMode()
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
	_, err := c.outStream.Write(p)
	if err != nil {
		log.Panic(err)
	}
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

func (c *ConsoleUI) GetCursor() (int, int) {
	c.Print(getCursorPosition)

	var bytes = make([]byte, 7)
	_, err := os.Stdin.Read(bytes)
	if err != nil {
		panic(err)
	}

	x, y, err := c.parsePosition(bytes)
	if err != nil {
		// TODO: remake
		log.Panic(err)
	}
	return x, y
}

func (c *ConsoleUI) parsePosition(bytes []byte) (int, int, error) {
	matches := c.cursorPositionRegex.FindStringSubmatch(string(bytes))
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
