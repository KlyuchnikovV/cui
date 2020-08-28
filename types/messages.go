package types

import (
	"log"
	"os"

	"github.com/KlyuchnikovV/termin/keys"
)

type Message interface {
	Exec(Widget)
}

type ResizeMsg struct {
	x, y, w, h *int
}

func NewResizeMsg(x, y, w, h *int) *ResizeMsg {
	if x != nil && *x < 1 {
		*x = 1
	}
	if y != nil && *y < 1 {
		*y = 1
	}
	if w != nil && *w < 1 {
		*w = 1
	}
	if h != nil && *h < 1 {
		*h = 1
	}
	return &ResizeMsg{
		x: x,
		y: y,
		w: w,
		h: h,
	}
}

func (r *ResizeMsg) Exec(w Widget) {
	log.Printf("%#v", *r.w)
	if r.x != nil {
		w.SetOption("x", *r.x)
	}
	if r.y != nil {
		w.SetOption("y", *r.y)
	}
	if r.w != nil {
		w.SetOption("w", *r.w)
	}
	if r.h != nil {
		w.SetOption("h", *r.h)
	}

	w.ClearScreen()
	switch w.GetIntOption("sizePolicy") {
	case 0:
		if err := w.DrawRectangle(w.X(), w.Y(), w.W(), w.H(), FullBlock); err != nil {
			w.SendError(err)
			return
		}
	}
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
	w.SetOption("signal", s.signal)
}

type KeyboardMsg struct {
	r keys.KeyboardKey
}

func NewKeyboardMsg(r keys.KeyboardKey) *KeyboardMsg {
	return &KeyboardMsg{
		r: r,
	}
}

func (k *KeyboardMsg) Exec(w Widget) {
	w.SetOption("rune", k.r)
}
