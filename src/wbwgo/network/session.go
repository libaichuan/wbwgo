package network

import (
	"bufio"
	"encoding/binary"
	"io"
	//"log"
	"net"
	"sync"

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

	wait_group sync.WaitGroup

	OnClose func()

	needPostSend bool
}

const (
	PackageHeaderSize = 4 // MsgID(uint16) + MsgLen(uint16)
	MaxPacketSize     = 2048
)

func (self *Session) WeakUpRecvLoop() {
	self.packet_list.AddPacket(&Packet{})
}

func (self *Session) Close() {
	if self.is_close {
		return
	}

	self.is_close = true

	self.conn.Close()

	//log.Println("session close\n")
}

func (self *Session) Send(data interface{}) {
	msg := data.(proto.Message)

	data_arr, err := proto.Marshal(msg)
	if err != nil {
		//log.Println("Send Error Msg %s", reflect.TypeOf(msg).String())
		return
	}

	cur_id := self.dispatcher.GetRefID(reflect.TypeOf(msg))
	if cur_id == 0 {
		return

	}

	p := &Packet{}
	p.data = data_arr

	p.msg_id = cur_id

	self.Write(p)
	//self.packetlist.AddPacket(p)
}

func (self *Session) Write(p *Packet) bool {
	msg_len := uint16(len(p.data))

	msg_buf := make([]byte, msg_len+4)

	binary.LittleEndian.PutUint16(msg_buf, p.msg_id)
	binary.LittleEndian.PutUint16(msg_buf[2:], msg_len)

	copy(msg_buf[4:], p.data)

	if _, err := self.conn.Write(msg_buf); err != nil {
		return false
	}
	return true
}

func (self *Session) RecvLoop() {
	defer func() {
		//log.Println("RecvLoop Over")
	}()

	for {
		if self.is_close {
			break
		}

		if _, err := io.ReadFull(self.head_reader, self.head_buf); err != nil {
			self.Close()
			//log.Println("%s", err)
			break
		}

		new_pack := &Packet{}

		new_pack.msg_id = binary.LittleEndian.Uint16(self.head_buf)
		new_pack.msg_len = binary.LittleEndian.Uint16(self.head_buf[2:])

		//log.Printf("msg_id:%d,msg_len:%d", new_pack.msg_id, new_pack.msg_len)

		if new_pack.msg_len > MaxPacketSize {
			//log.Println("more than maxpackage")
			break
		}

		new_pack.data = make([]byte, new_pack.msg_len)

		if _, err := io.ReadFull(self.head_reader, new_pack.data); err != nil {
			//log.Println("%s", err)
			self.Close()
			break
		}
		//log.Println("收到data:", new_pack.data)

		self.event_loop.AddInLoop(self.dispatcher, NewEventData(self, new_pack))

	}

	//	if self.needPostSend {
	//		self.WeakUpRecvLoop()
	//	}

	self.wait_group.Done()
}

func (self *Session) SendLoop() {
	defer func() {
		//log.Println("SendLoop Over")
	}()

	var write_list []*Packet

	for {
		if self.is_close {
			break
		}

		will_exit := false

		write_list = write_list[0:0]

		temp_list := self.packet_list.BeginList()
		for _, v := range temp_list {
			if v.msg_id == 0 {
				will_exit = true
			} else {
				write_list = append(write_list, v)
			}
		}
		self.packet_list.EndList()

		for _, v := range write_list {
			if self.Write(v) {
				break
			}
		}

		if will_exit {
			break
		}
	}

	self.Close()

	self.needPostSend = false

	self.wait_group.Done()
}

func NewSession(conn net.Conn, dispatcher *MsgDispatcher, event_loop *EventLoop) *Session {
	se := &Session{
		conn:         conn,
		head_reader:  bufio.NewReader(conn),
		head_buf:     make([]byte, PackageHeaderSize),
		is_close:     false,
		event_loop:   event_loop,
		dispatcher:   dispatcher,
		packet_list:  NewPacketList(),
		needPostSend: true,
	}

	se.wait_group.Add(1)

	go func() {
		se.wait_group.Wait()
		se.OnClose()
	}()

	go se.RecvLoop()

	//go se.SendLoop()

	return se
}
