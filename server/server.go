package server

import (
	"context"
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

	*graphics.Graphics
}

func New(ctx context.Context, widgets ...types.Widget) *Server {
	return &Server{
		ctx:      ctx,
		widgets:  widgets,
		Graphics: graphics.New(),
		ch:       make(chan chan_utils.Message, 10),
	}
}

func (s *Server) RegisterWidgets(widgets ...types.Widget) {
	s.widgets = append(s.widgets, widgets...)
}

func (s *Server) StartRendering(async bool) {
	if s.cancel != nil {
		return
	}

	render, cancel := chan_utils.NewListener(s.ctx, s.ch, s.onRenderRequest, func(err error) {
		log.Panic(err)
	})

	s.cancel = cancel

	// Listening to window resize
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go s.listenSysCalls(ch)

	if async {
		go render()
	} else {
		render()
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
		panic("!ok")
	}
	log.Printf("INFO: got to render %#v", string(msg))
	_, err := s.Write(msg)
	if err != nil {
		panic(err)
	}
}

func (s *Server) GetRenderChan() chan chan_utils.Message {
	return s.ch
}

func (s *Server) listenSysCalls(ch chan os.Signal) {
	for {
		select {
		case request, ok := <-ch:
			if !ok {
				log.Panic("sig channel was closed")
			}
			for _, listener := range s.widgets {
				listener.Render(request)
			}
		case <-s.ctx.Done():
			return
		}
	}
}
