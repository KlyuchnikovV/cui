package widgets

//
import (
	"log"

	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/server"
	"github.com/KlyuchnikovV/cui/types"
)

type Textarea struct {
	baseElement
	x, y int
	text string
}

func NewTextarea(c *cui.ConsoleUI, text string) *Textarea {
	t := &Textarea{
		baseElement: *newBaseElement(c, nil),
		x:           0,
		y:           0,
		text:        text,
	}
	// TODO: optimize subscription mechanism
	c.SubscribeWidget(server.KeyboardChan, t)
	c.SubscribeWidget(server.ResizeChan, t)
	return t
}

func (t *Textarea) Render(msg types.Message) {
	msg.Exec(t)

	switch msg.(type) {
	case *types.KeyboardMsg:
		value, ok := t.GetOption("rune").(rune)
		if !ok {
			log.Println("WARN: rune not set")
			return
		}
		log.Printf("textarea got rune %s\n", string(value))
		t.processRune(value)
	case *types.ResizeMsg:
		x, y, w, h := t.GetIntOption("x"), t.GetIntOption("y"), t.GetIntOption("w"), t.GetIntOption("h")
		log.Printf("%s draw at %d %d %d %d\n", t.text, x, y, h, w)
		t.PrintAt(x+2, y+1, t.text, true)
	}
}

func (t *Textarea) processRune(r rune) bool {
	x, y := t.GetIntOption("x"), t.GetIntOption("y")
	switch {
	case r == 10:
		t.x += 1
		t.y = 0
		t.text += "\n"
		// TODO: make getting relative cursor position more easy
		t.SetCursor(x+t.x+1, y+t.y+1)
	case r >= 0 && r <= 31:
	case r > 31 && r <= 127:
		t.y += 1
		t.text += string(r)
		log.Printf("INFO textarea's text %s\n", t.text)
		t.PrintAt(x+t.x+1, y+t.y, string(r), false)
	}
	return false
}
