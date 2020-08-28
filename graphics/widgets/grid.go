package widgets

import (
	"log"

	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/server"
	"github.com/KlyuchnikovV/cui/types"
)

type Grid struct {
	baseElement
	children []gridRow
}

type gridRow []types.Widget

func newGridRow(c *cui.ConsoleUI, widgets ...func(*cui.ConsoleUI) types.Widget) *gridRow {
	var row = make(gridRow, len(widgets))
	for i, w := range widgets {
		row[i] = w(c)
	}
	return &row
}

type GridWidgetLayout [][]func(*cui.ConsoleUI) types.Widget

func NewGrid(widgets GridWidgetLayout) func(c *cui.ConsoleUI) types.Widget {
	return func(c *cui.ConsoleUI) types.Widget {
		var rows = make([]gridRow, len(widgets))
		for i, row := range widgets {
			rows[i] = *newGridRow(c, row...)
		}
		g := &Grid{
			baseElement: *newBaseElement(c, nil),
			children:    rows,
		}
		c.SubscribeWidget(server.ResizeChan, g)
		return g
	}
}

func (g *Grid) Render(msg types.Message) {
	log.Print("grid: exec")
	msg.Exec(g)

	childH := g.H() / len(g.children)
	log.Printf("grid: calc childH: %d", childH)

	x := g.X()

	for i, row := range g.children {
		// For division by two fix
		if i == len(g.children)-1 && (i+1)*childH < g.H() {
			x -= 1
			childH += 1
		}

		childW := g.W() / len(row)
		log.Printf("grid: calc i: %d childW: %d", i, childW)

		for j, child := range row {
			if j == len(row)-1 && (j+1)*childW < g.W() {
				childW += 1
			}
			newX := x + i*childH + 1
			newY := g.Y() + j*childW + 1
			child.Render(types.NewResizeMsg(&newX, &newY, &childW, &childH))
		}
	}
}
