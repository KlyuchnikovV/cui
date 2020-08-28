package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlyuchnikovV/cui/graphics"
	"github.com/KlyuchnikovV/cui/low_level/terminal"
	"github.com/KlyuchnikovV/cui/types"
	"github.com/KlyuchnikovV/termin"
)

type Server struct {
	ctx     context.Context
	cancel  context.CancelFunc
	ch      chan types.Message
	input   *termin.Termin
	Widgets map[ChanEnum][]types.Widget

	types.ErrorChannel
	*graphics.Graphics
}

func New(ctx context.Context, widgets map[ChanEnum][]types.Widget) *Server {
	if widgets == nil {
		widgets = make(map[ChanEnum][]types.Widget)
	}
	return &Server{
		ctx:          ctx,
		Widgets:      widgets,
		Graphics:     graphics.New(),
		ch:           make(chan types.Message, 1),
		ErrorChannel: types.NewErrorChannel(1),
		input:        termin.New(),
	}
}

func (s *Server) SubscribeWidget(key ChanEnum, widget types.Widget) {
	s.Widgets[key] = append(s.Widgets[key], widget)
}

func (s *Server) StartRendering(async bool) {
	if s.cancel != nil {
		return
	}
	s.ctx, s.cancel = context.WithCancel(s.ctx)

	// Listening to keys
	s.input.StartReading(true)

	if async {
		go s.update()
	} else {
		s.update()
	}
}

func (s *Server) Stop() {
	if s.cancel == nil {
		return
	}

	s.cancel()
	s.cancel = nil
}

func (s *Server) GetRenderChan() chan types.Message {
	return s.ch
}

func (s *Server) update() {
	defer func() {
		if e := recover(); e != nil {
			err, ok := e.(error)
			if !ok {
			err = fmt.Errorf("%#v", e)
			}
			s.SendError(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	// Listening to window resize
	signal.Notify(ch, syscall.SIGWINCH)

	// TODO: Initial update
	ch <- syscall.SIGWINCH

	for {
		select {
		case signal, ok := <-ch:
			if !ok {
				s.SendError(fmt.Errorf("signal channel was unexpectedly closed"))
				return
			}

			if signal != syscall.SIGWINCH {
				s.SendError(fmt.Errorf("signal channel was unexpectedly closed"))
				continue
			}

			w, h := terminal.GetTerminalSize()
			s.broadcast(types.NewResizeMsg(nil, nil, &w, &h), s.Widgets[ResizeChan]...)
		case msg, ok := <-s.input.GetChan():
			if !ok {
				s.SendError(fmt.Errorf("input channel was unexpectedly closed"))
				return
			}
			s.broadcast(types.NewKeyboardMsg(msg), s.Widgets[KeyboardChan]...)
		case msg, ok := <-s.ch:
			if !ok {
				s.SendError(fmt.Errorf("resize channel was unexpectedly closed"))
				return
			}
			s.broadcast(msg, s.Widgets[ResizeChan]...)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Server) broadcast(msg types.Message, widgets ...types.Widget) {
	for _, widget := range widgets {
		widget.Render(msg)
	}
}
