package widgets

import (
	"fmt"
	"os"
	"syscall"

	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/low_level/terminal"
	"github.com/KlyuchnikovV/cui/types"
)

type Body struct {
	baseElement
}

func NewBody(widgetGenerators ...WidgetGenerator) WidgetGenerator {
	return func(c *cui.ConsoleUI) types.Widget {
		b := &Body{
			baseElement: *newBaseElement(c, nil, widgetGenerators...),
		}
		return b
	}
}

func (b *Body) Render(msg types.Message) {
	b.SavePosition()
	msg.Exec(b)

	x, y := b.X(), b.Y()
	w, h := terminal.GetTerminalSize()

	signal, ok := b.GetOption("signal").(os.Signal)
	if !ok {
		b.SendError(fmt.Errorf("BODY: signal not of type %T", signal))
		return
	}
	switch signal {
	case syscall.SIGWINCH:
		b.SetOption("w", w)
		b.SetOption("h", h)

		for _, child := range b.children {
			child.Render(types.NewResizeMsg(&x, &y, &w, &h))
		}
	}
	b.RestorePosition()
}
