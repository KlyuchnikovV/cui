package server

type ChanEnum interface {
	getValue() int
}

type chanEnum int

func (c chanEnum) getValue() int {
	return int(c)
}

const (
	KeyboardChan chanEnum = iota
	ResizeChan
)
