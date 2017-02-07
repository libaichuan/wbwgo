package network

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Session struct {
	id int32

	conn net.Conn

	event_loop *EventLoop

	dispatcher *MsgDispatcher

	head_reader *bufio.Reader

	head_buf []byte

	is_close bool
}

const (
	PackageHeaderSize = 4 // MsgID(uint16) + MsgLen(uint16)
	MaxPacketSize     = 2048
)

func (self *Session) Close() {
	if self.is_close {
		log.Fatalln("close double")
		return
	}

	self.is_close = true

	self.conn.Close()
}

func (self *Session) RecvLoop() {

	for {
		if self.is_close {
			break
		}

		if _, err := io.ReadFull(self.head_reader, self.head_buf); err != nil {
			log.Fatalln("%s", err)
			return
		}

		new_pack := &Packet{}

		new_pack.msg_id = binary.LittleEndian.Uint16(self.head_buf)
		new_pack.msg_len = binary.LittleEndian.Uint16(self.head_buf[2:])

		if new_pack.msg_len > MaxPacketSize {
			log.Fatalln("more than maxpackage")
			return
		}

		new_pack.data = make([]byte, new_pack.msg_len)

		if _, err := io.ReadFull(self.conn, new_pack.data); err != nil {
			log.Fatalln("%s", err)
			return
		}

		/*to do
		发给业务线程的队列，然后反序列.开始想在这里反序列化，然后直接发proto类到业务线程的队列。
		但是这里要转一次基类，到时候又不能确定当然proto的回调。要是通过ID再在相应的业务里面转一次，转的次数太多，而且没泛型。
		所以这里就不反序列化了，在回调里面自己转
		*/
		self.event_loop.AddInLoop(self.dispatcher, new_pack)
	}
}

func (self *Session) SendLoop() {

}

func NewSession(conn net.Conn, dispatcher *MsgDispatcher, event_loop *EventLoop) *Session {
	se := &Session{
		conn:        conn,
		head_reader: bufio.NewReader(conn),
		head_buf:    make([]byte, PackageHeaderSize),
		is_close:    false,
		event_loop:  event_loop,
		dispatcher:  dispatcher,
	}

	go se.RecvLoop()

	go se.SendLoop()

	return se
}
