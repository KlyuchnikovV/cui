package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlyuchnikovV/cui/graphics"
	"github.com/KlyuchnikovV/cui/types"
)

type Server struct {
	ctx     context.Context
	cancel  context.CancelFunc
	ch      chan types.Message
	widgets []types.Widget

	types.ErrorChannel
	*graphics.Graphics
}

func New(ctx context.Context, widgets ...types.Widget) *Server {
	return &Server{
		ctx:          ctx,
		widgets:      widgets,
		Graphics:     graphics.New(),
		ch:           make(chan types.Message, 1),
		ErrorChannel: types.NewErrorChannel(1),
	}
}

func (s *Server) RegisterWidgets(widgets ...types.Widget) {
	s.widgets = append(s.widgets, widgets...)
}

func (s *Server) StartRendering(async bool) {
	if s.cancel != nil {
		return
	}
	s.ctx, s.cancel = context.WithCancel(s.ctx)

	ch := make(chan os.Signal, 1)
	// Listening to window resize
	signal.Notify(ch, syscall.SIGWINCH)
	go s.redirectSignals(ch)

	if async {
		go s.update()
	} else {
		s.update()
	}
}

func (s *Server) update() {
	defer func() {
		if e := recover(); e != nil {
			s.SendError(e.(error))
		}
	}()

	for _, listener := range s.widgets {
		listener.Render(types.NewSignalMsg(syscall.SIGWINCH))
	}

	for {
		select {
		case msg, ok := <-s.ch:
			if !ok {
				s.SendError(fmt.Errorf("update channel was unexpectedly closed"))
				return
			}

			log.Printf("TRACE: got to update %#v", msg)
			for _, widget := range s.widgets {
				widget.Render(msg)
			}
		case <-s.ctx.Done():
			return
		}
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

func (s *Server) redirectSignals(ch chan os.Signal) {
	defer func() {
		if e := recover(); e != nil {
			s.SendError(e.(error))
		}
	}()

	for {
		select {
		case signal, ok := <-ch:
			if !ok {
				s.SendError(fmt.Errorf("signal channel was unexpectedly closed"))
				return
			}

			log.Printf("TRACE: got signal %#v", signal)
			s.SetCursor(0, 0)
			s.ClearScreen(graphics.ClearAfterCursor)
			for _, listener := range s.widgets {
				listener.Render(types.NewSignalMsg(signal))
			}
		case <-s.ctx.Done():
			return
		}
	}
}
