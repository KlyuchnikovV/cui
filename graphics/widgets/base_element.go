package widgets

import (
	"strings"

	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/types"
)

type WidgetGenerator func(c *cui.ConsoleUI) types.Widget

type baseElement struct {
	*cui.ConsoleUI
	options  map[string]interface{}
	children []types.Widget
}

func newBaseElement(c *cui.ConsoleUI, options map[string]interface{}, children ...WidgetGenerator) *baseElement {
	var widgets = make([]types.Widget, len(children))
	for i, widgetGenerator := range children {
		widgets[i] = widgetGenerator(c)
	}

	if options == nil {
		options = make(map[string]interface{})
	}
	return &baseElement{
		ConsoleUI: c,
		options:   options,
		children:  widgets,
	}
}

func (b *baseElement) SetOption(key string, value interface{}) {
	b.options[key] = value
}

func (b *baseElement) SetOptions(options map[string]interface{}) {
	for key, value := range options {
		b.SetOption(key, value)
	}
}

func (b *baseElement) GetOption(s string) interface{} {
	return b.options[s]
}

func (b *baseElement) GetOptions() map[string]interface{} {
	return b.options
}

func (b *baseElement) GetIntOption(s string) int {
	opt := b.options[s]
	if opt == nil {
		return 0
	}
	result, ok := opt.(int)
	if !ok {
		return 0
	}
	return result
}

func (b *baseElement) ClearScreen() {
	b.SavePosition()
	var replaceString = strings.Repeat(" ", b.W()-2)
	for i := b.X() + 1; i < b.H(); i++ {
		b.PrintAt(i, b.Y()+1, replaceString, false)
	}
	b.RestorePosition()
}

func (b *baseElement) X() int {
	x := b.GetIntOption("x")
	if x < 1 {
		x = 1
	}
	return x
}

func (b *baseElement) Y() int {
	y := b.GetIntOption("y")
	if y < 1 {
		y = 1
	}
	return y
}

func (b *baseElement) W() int {
	w := b.GetIntOption("w")
	if w < 1 {
		w = 1
	}
	return w
}

func (b *baseElement) H() int {
	h := b.GetIntOption("h")
	if h < 1 {
		h = 1
	}
	return h
}
