package widgets

import (
	"os"
	"syscall"

	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/low_level/terminal"
	"github.com/KlyuchnikovV/cui/types"
)

type Body struct {
	baseElement
}

func NewBody(c *cui.ConsoleUI, children ...types.Widget) *Body {
	return &Body{
		baseElement: *newBaseElement(c, nil, children...),
	}
}

func (b *Body) Render(msg types.Message) {
	msg.Exec(b)
	switch b.GetOption("signal").(os.Signal) {
	case syscall.SIGWINCH:
		b.options["w"], b.options["h"] = terminal.GetTerminalSize()
		for _, child := range b.children {
			child.Render(types.NewResizeMsg(b.GetIntOption("x"), b.GetIntOption("y"), b.GetIntOption("w"), b.GetIntOption("h")))
		}
	}
}
