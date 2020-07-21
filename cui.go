package cui

import (
	"github.com/KlyuchnikovV/cui/server"
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

func New(ctx context.Context, enableRaw bool, widgets ...types.Widget) (*ConsoleUI, error) {
	if enableRaw {
		log.Print("entering raw")
		raw_mode.EnableRawMode()
	}

	return &ConsoleUI{
		Server: server.New(ctx, widgets...)
	}, nil
}

func (c *ConsoleUI) DisableRawMode() {
	raw_mode.DisableRawMode()
}
