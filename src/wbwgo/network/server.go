package network

import (
	"log"
	"net"
)

type Server struct {
	*SessionManager

	listener net.Listener

	running bool
}

func NewServer() *Server {
	return &Server{
		SessionManager: NewSessionManager(),
		running:        false,
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

			conn.LocalAddr().String()
		}
	}()

}
