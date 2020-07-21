package types

import (
	"os"

	"github.com/KlyuchnikovV/chan_utils"
)

type Message []byte

func (r Message) GetMessage() chan_utils.Message {
	return r
}

func (r Message) String() string {
	return string(r)
}

func (r Message) Signal() {}

type Widget interface {
	Render(Message)
	ProcessSystemSignal(os.Signal)
}

type ConsoleStream interface {
	getConsoleStream() ConsoleStream
	Print(string) error
	Write([]byte) (int, error)
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

func (c *consoleStream) Write(p []byte) (int, error) {
	return c.out.Write(p)
}

func (c *consoleStream) Print(message string) error {
	_, err := c.out.Write([]byte(message))
	return err
}

func (c *consoleStream) getConsoleStream() ConsoleStream {
	return c
}
