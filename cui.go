package cui

import (
	"context"
	"log"

	"github.com/KlyuchnikovV/cui/server"
	"github.com/KlyuchnikovV/cui/types"
	"github.com/KlyuchnikovV/termin/low_level/raw_mode"
)

type ConsoleUI struct {
	server.Server
	body types.Widget
}

func New(ctx context.Context, enableRaw bool, widgets map[server.ChanEnum][]types.Widget) *ConsoleUI {
	if enableRaw {
		log.Print("entering raw")
		raw_mode.EnableRawMode()
	}

	return &ConsoleUI{
		Server: *server.New(ctx, widgets),
	}
}

func (c *ConsoleUI) DisableRawMode() {
	raw_mode.DisableRawMode()
}

func (c *ConsoleUI) Layout(generator func(*ConsoleUI) types.Widget) {
	c.body = generator(c)
}
