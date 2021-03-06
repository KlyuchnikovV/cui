package types

import (
	"os"
)

type Widget interface {
	Render(Message)
	SetOptions(map[string]interface{})
	GetOption(string) interface{}
	GetIntOption(string) int
	GetOptions() map[string]interface{}
	SendError(error)
	DrawRectangle(int, int, int, int, rune) error
	// PrintAt(int, int, string, bool)
}

type ConsoleStream interface {
	getConsoleStream() ConsoleStream
	Print(string) error
	Write([]byte) (int, error)
	Stream() *os.File
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

func (c *consoleStream) Stream() *os.File {
	return c.out
}
