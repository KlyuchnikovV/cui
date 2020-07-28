package widgets

import (
	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/types"
)

type Grid struct {
	baseElement
	children [][]types.Widget
}

func NewGrid(c *cui.ConsoleUI, children [][]types.Widget) *Grid {
	return &Grid{
		baseElement: *newBaseElement(c, nil, nil),
		children:    children,
	}
}

func (g *Grid) Render(msg types.Message) {
	msg.Exec(g)
	x, y, w, h := g.GetIntOption("x"), g.GetIntOption("y"), g.GetIntOption("w"), g.GetIntOption("h")
	childH := h / len(g.children)
	for i, row := range g.children {
		// For division by two fix
		if i == len(g.children)-1 && (i+1)*childH < h {
			x -= 1
			childH += 1
		}
		childW := w / len(row)
		for j, child := range row {
			if j == len(row)-1 && (j+1)*childW < w {
				childW += 1
			}
			child.Render(types.NewResizeMsg(x+i*childH, y+j*childW, childW, childH))
		}
	}
}
