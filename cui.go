package cui

import (
	"github.com/KlyuchnikovV/cui/graphics"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/KlyuchnikovV/chan_utils"
	"github.com/KlyuchnikovV/cui/low_level/raw_mode"
)

type ConsoleUI struct {
	server.Server
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

	return &ConsoleUI{
		ctx:                 ctx,
		Graphics:            graphics.New()
		renderChan:          make(chan chan_utils.Message, 10),
		cursorPositionRegex: r,
	}, nil
}

func (c *ConsoleUI) DisableRawMode() {
	raw_mode.DisableRawMode()
}
