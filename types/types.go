package types

import (
	"os"

	"github.com/KlyuchnikovV/chan_utils"
)

type RenderRequest []byte

func (r *RenderRequest) GetMessage() chan_utils.Message {
	return r
}

type Widget interface {
	Render(RenderRequest)
}

type ConsoleStream interface {
	getConsoleStream() ConsoleStream
	Print(string) error
}

type consoleStream struct {
	out *os.File
}

func NewConsoleStream() *consoleStream {
	result := consoleStream{
		out: os.Stdout,
	}
	return &result
}

func (c *consoleStream) Print(message string) error {
	_, err := c.out.Write([]byte(message))
	return err
}

func (c *consoleStream) getConsoleStream() ConsoleStream {
	return c
}
