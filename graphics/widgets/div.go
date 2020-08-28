package widgets

import (
	"log"

	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/types"
)

type Div struct {
	baseElement
	text string
}

func NewDiv(text string) func(c *cui.ConsoleUI) types.Widget {
	return func(c *cui.ConsoleUI) types.Widget {
		return &Div{
			baseElement: *newBaseElement(c, nil),
			text:        text,
		}
	}
}

func (d *Div) Render(msg types.Message) {
	log.Print("div: exec")
	msg.Exec(d)
	d.PrintAt(d.X()+d.H()/2, d.Y()+d.W()/2, d.text, true)
	log.Printf("div: print %s at %d %d", d.text, d.X()+d.H()/2, d.Y()+d.W()/2)
}
