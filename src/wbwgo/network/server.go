package network

import (
	"log"
	"net"
)

type Server struct {
	ses_manager *SessionManager

	event_loop *EventLoop

	dispatcher *MsgDispatcher

	listener net.Listener

	running bool
}

func NewServer(event_loop *EventLoop) *Server {
	return &Server{
		ses_manager: NewSessionManager(),
		event_loop:  event_loop,
		dispatcher:  NewMsgDispatcher(),
		running:     false,
	}
}

func (s *Server) Init(conn_type string, addr string) {
	ln, err := net.Listen(conn_type, addr)

	if err != nil {
		log.Fatalf("listen error %s", err.Error())
		return
	}

	s.listener = ln
	s.running = true

	go func() {
		for s.running {
			conn, err := ln.Accept()

			if err != nil {
				log.Fatalf("accpet error %s", err.Error())
				continue
			}

			se := NewSession(conn)

			s.ses_manager.AddSession(se)

			//todo session关闭

			log.Printf("accept new conn id:%d", se.id)
		}
	}()
}
