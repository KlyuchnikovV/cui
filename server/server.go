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

	"github.com/KlyuchnikovV/chan_utils"
)

type Server struct {
	ctx     context.Context
	cancel  context.CancelFunc
	ch      chan chan_utils.Message
	widgets []types.Widget

	types.ErrorChannel
	*graphics.Graphics
}

func New(ctx context.Context, widgets ...types.Widget) (*Server, error) {
	g, err := graphics.New()
	if err != nil {
		return nil, err
	}
	return &Server{
		ctx:          ctx,
		widgets:      widgets,
		Graphics:     g,
		ch:           make(chan chan_utils.Message, 1),
		ErrorChannel: types.NewErrorChannel(1),
	}, nil
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
		go s.render()
	} else {
		s.render()
	}
}

func (s *Server) render() {
	defer func() {
		if e := recover(); e != nil {
			s.SendError(e.(error))
		}
	}()

	for {
		select {
		case msg, ok := <-s.ch:
			if !ok {
				s.SendError(fmt.Errorf("channel was closed"))
				return
			}

			s.onRenderRequest(msg)
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

func (s *Server) onRenderRequest(data chan_utils.Message) {
	msg, ok := data.(types.RenderRequest)
	if !ok {
		s.SendError(fmt.Errorf("data wasn't of type \"%T\"", msg))
		return
	}
	log.Printf("INFO: got to render %#v", string(msg))
	if _, err := s.Write(msg); err != nil {
		s.SendError(err)
	}
}

func (s *Server) GetRenderChan() chan chan_utils.Message {
	return s.ch
}

func (s *Server) redirectSignals(ch chan os.Signal) {
	for {
		select {
		case signal, ok := <-ch:
			if !ok {
				s.SendError(fmt.Errorf("signal channel was unexpectedly closed"))
				return
			}
			for _, listener := range s.widgets {
				listener.ProcessSystemSignal(signal)
			}
		case <-s.ctx.Done():
			return
		}
	}
}
