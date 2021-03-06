package network

import (
	//"log"
	"reflect"

	"github.com/golang/protobuf/proto"
)

/*
负责事件的注册和分发
*/

type MsgDispatcher struct {
	msgs map[uint16]func(interface{})

	ref_types map[reflect.Type]uint16
}

func (self *MsgDispatcher) RegisterMsg(id uint16, f func(interface{})) {
	if _, ok := self.msgs[id]; ok {
		//log.Fatalln("MsgDispatcher::RegisterMsg same id")
		return
	}

	self.msgs[id] = f
}

func (self *MsgDispatcher) RegisterRefType(t reflect.Type, id uint16) {
	if _, ok := self.ref_types[t]; ok {
		//log.Fatalln("MsgDispatcher::RegisterRefType same reflect.Type")
		return
	}

	self.ref_types[t] = id
}

func (self *MsgDispatcher) GetRefID(t reflect.Type) uint16 {
	if v, ok := self.ref_types[t]; ok {
		return v
	}

	return 0
}

func (self *MsgDispatcher) OnMessage(data interface{}) {
	//log.Println("OnMessage ...................")
	if v, ok := data.(*EventData); ok {
		if f, ok := self.msgs[v.p.msg_id]; ok {
			f(data)
		} else {
			//log.Println("MsgDispatcher::OnMessage not register message")
		}
	} else {
		//log.Println("MsgDispatcher::OnMessage not EventData")
	}

}

func NewMsgDispatcher() *MsgDispatcher {
	self := &MsgDispatcher{
		msgs:      make(map[uint16]func(interface{})),
		ref_types: make(map[reflect.Type]uint16),
	}

	return self
}

func RegisterMessage(dispatcher *MsgDispatcher, id uint16, m_msg proto.Message, callback func(ses *Session, f interface{})) {
	ref_type := reflect.TypeOf(m_msg)

	dispatcher.RegisterRefType(ref_type, id)

	dispatcher.RegisterMsg(id, func(data interface{}) {

		if ev, ok := data.(*EventData); ok {
			m_msg := reflect.New(ref_type.Elem()).Interface()
			proto.Unmarshal(ev.p.data, m_msg.(proto.Message))

			callback(ev.ses, m_msg)
		} else {
			//log.Printf("dispatcher callback to *EventData,data type is %s", reflect.TypeOf(data).String())
		}

	})
}
