package types

import (
	// "log"
	"os"

	"github.com/KlyuchnikovV/termin/keys"
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
	// TODO: terminal count symbols from 1
	if r.x < 1 {
		r.x = 1
	}
	if r.y < 1 {
		r.y = 1
	}
	w.SetOptions(map[string]interface{}{
		"x": r.x,
		"y": r.y,
		"w": r.w,
		"h": r.h,
	})
	switch w.GetIntOption("sizePolicy") {
	case 0:
		if err := w.DrawRectangle(r.x, r.y, r.w, r.h, FullBlock); err != nil {
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
	w.SetOptions(map[string]interface{}{
		"signal": s.signal,
	})
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
	w.SetOptions(map[string]interface{}{
		"rune": k.r,
	})
}
