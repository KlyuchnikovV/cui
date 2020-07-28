package types

import (
	"os"
)

type Message interface {
	Exec(Widget)
}

type ResizeMsg struct {
	x, y int
	w, h int
}

func NewResizeMsg(x, y, w, h int) *ResizeMsg {
	return &ResizeMsg{
		x: x,
		y: y,
		w: w,
		h: h,
	}
}

func (r *ResizeMsg) Exec(w Widget) {
	w.SetOptions(map[string]interface{}{
		"x": r.x,
		"y": r.y,
		"w": r.w,
		"h": r.h,
	})
}

type SignalMsg struct {
	signal os.Signal
}

func NewSignalMsg(signal os.Signal) *SignalMsg {
	return &SignalMsg{
		signal: signal,
	}
}

func (s *SignalMsg) Exec(w Widget) {
	w.SetOptions(map[string]interface{}{
		"signal": s.signal,
	})
}

type Widget interface {
	Render(Message)
	// ProcessSystemSignal(os.Signal)
	SetOptions(opts map[string]interface{})
	GetOption(s string) interface{}
	GetOptions() map[string]interface{}
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
