package cui

import (
	"context"
	"log"

	"github.com/KlyuchnikovV/cui/server"
	"github.com/KlyuchnikovV/cui/types"

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

	s, err := server.New(ctx, widgets...)
	if err != nil {
		return nil, err
	}

	return &ConsoleUI{
		Server: *s,
	}, nil
}

func (c *ConsoleUI) DisableRawMode() {
	raw_mode.DisableRawMode()
}
