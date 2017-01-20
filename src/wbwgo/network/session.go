package network

import (
	"net"
)

type Session struct {
	id int32

	conn net.Conn
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn: conn,
	}
}
