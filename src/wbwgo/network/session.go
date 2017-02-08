package network

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"

	"reflect"

	"github.com/golang/protobuf/proto"
)

type Session struct {
	id int32

	conn net.Conn

	event_loop *EventLoop

	dispatcher *MsgDispatcher

	head_reader *bufio.Reader

	head_buf []byte

	is_close bool

	packet_list *PacketList
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
	log.Println("3333333333333333")

	self.is_close = true

	self.conn.Close()
}

func (self *Session) Send(data interface{}) {
	msg := data.(proto.Message)

	data_arr, err := proto.Marshal(msg)
	if err != nil {
		log.Println("Send Error Msg %s", reflect.TypeOf(msg).String())
		return
	}

	cur_id := self.dispatcher.GetRefID(reflect.TypeOf(msg))
	if cur_id == 0 {
		return

	}

	p := &Packet{}
	p.data = data_arr

	p.msg_id = cur_id

	self.packet_list.AddPacket(p)
}

func (self *Session) Write(p *Packet) {
	msg_len := uint16(len(p.data))

	msg_buf := make([]byte, msg_len+4)

	binary.LittleEndian.PutUint16(msg_buf, p.msg_id)
	binary.LittleEndian.PutUint16(msg_buf[2:], msg_len)

	copy(msg_buf[4:], p.data)

	if _, err := self.conn.Write(msg_buf); err != nil {
		log.Println(err)
	}
}

func (self *Session) RecvLoop() {

	for {
		if self.is_close {
			break
		}

		if _, err := io.ReadFull(self.conn, self.head_buf); err != nil {
			log.Println("111111111111")
			log.Println("%s", err)
			return
		}

		new_pack := &Packet{}

		new_pack.msg_id = binary.LittleEndian.Uint16(self.head_buf)
		new_pack.msg_len = binary.LittleEndian.Uint16(self.head_buf[2:])

		log.Printf("msg_id:%d,msg_len:%d", new_pack.msg_id, new_pack.msg_len)

		if new_pack.msg_len > MaxPacketSize {
			log.Println("more than maxpackage")
			return
		}

		new_pack.data = make([]byte, new_pack.msg_len)

		if _, err := io.ReadFull(self.conn, new_pack.data); err != nil {
			log.Println("2222222222")
			log.Println("%s", err)
			return
		}
		log.Println("收到data:", new_pack.data)

		self.event_loop.AddInLoop(self.dispatcher, NewEventData(self, new_pack))
	}
}

func (self *Session) SendLoop() {
	var write_list []*Packet

	for {
		if self.is_close {
			break
		}

		write_list = write_list[0:0]

		temp_list := self.packet_list.BeginList()
		for _, v := range temp_list {
			write_list = append(write_list, v)
		}
		self.packet_list.EndList()

		for _, v := range write_list {
			self.Write(v)
		}
	}
}

func NewSession(conn net.Conn, dispatcher *MsgDispatcher, event_loop *EventLoop) *Session {
	se := &Session{
		conn:        conn,
		head_reader: bufio.NewReader(conn),
		head_buf:    make([]byte, PackageHeaderSize),
		is_close:    false,
		event_loop:  event_loop,
		dispatcher:  dispatcher,
		packet_list: NewPacketList(),
	}

	go se.RecvLoop()

	go se.SendLoop()

	return se
}
