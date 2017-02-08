package network

import (
	"sync"
)

type PacketList struct {
	packet_list []*Packet
	list_guard  sync.Mutex
	list_cond   *sync.Cond
}

func (self *PacketList) AddPacket(p *Packet) {
	self.list_guard.Lock()
	defer self.list_guard.Unlock()

	self.packet_list = append(self.packet_list, p)

	self.list_cond.Signal()
}

func (self *PacketList) EndList() {
	self.packet_list = self.packet_list[0:0]

	self.list_guard.Unlock()
}

func (self *PacketList) BeginList() []*Packet {
	self.list_guard.Lock()

	for len(self.packet_list) == 0 {
		self.list_cond.Wait()
	}

	self.list_guard.Unlock()

	self.list_guard.Lock()

	return self.packet_list
}

func NewPacketList() *PacketList {
	pl := &PacketList{}
	pl.list_cond = sync.NewCond(&pl.list_guard)

	return pl
}
