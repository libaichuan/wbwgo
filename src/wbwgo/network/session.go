package network

import (
	"net"
	"sync/atomic"
)

var globalSessionId int = 0

type Session struct {
	id int

	conn net.Conn
}

func NewSession(conn net.Conn) *Session {
	self := &Session{
		conn:		conn,
		id:			atomic.AddInt32(&globalSessionId,1)
	}
}
