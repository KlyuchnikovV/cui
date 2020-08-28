package widgets

import (
	"strings"

	"github.com/KlyuchnikovV/cui"
	"github.com/KlyuchnikovV/cui/types"
)

type baseElement struct {
	*cui.ConsoleUI
	options  map[string]interface{}
	children []types.Widget
}

func newBaseElement(c *cui.ConsoleUI, options map[string]interface{}, children ...func(c *cui.ConsoleUI) types.Widget) *baseElement {
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
	var replaceString = strings.Repeat(" ", b.W()-3)
	for i := b.X() + 1; i < b.H(); i++ {
		b.PrintAt(i, b.Y()+1, replaceString, false)
	}
	b.RestorePosition()
}

func (b *baseElement) X() int {
	return b.GetIntOption("x")
}

func (b *baseElement) Y() int {
	return b.GetIntOption("y")
}

func (b *baseElement) W() int {
	return b.GetIntOption("w")
}

func (b *baseElement) H() int {
	return b.GetIntOption("h")
}
