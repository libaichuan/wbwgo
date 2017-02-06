package network

import (
	"log"
)

/*
负责事件的注册和分发
*/

type MsgDispatcher struct {
	msgs map[uint16]func(interface{})
}

func (self *MsgDispatcher) RegisterMsg(id uint16, f func(interface{})) {
	if _, ok := self.msgs[id]; ok {
		log.Fatalln("MsgDispatcher::RegisterMsg same id")
		return
	}

	self.msgs[id] = f
}

func (self *MsgDispatcher) OnMessage(data interface{}) {

}
