package main

import (
	"log"
	"wbwgo/msg"
	"wbwgo/network"

	"reflect"

	"github.com/golang/protobuf/proto"
)

var dis *network.MsgDispatcher

func RegisterMessage(id uint16, m_msg proto.Message) {
	dis.RegisterRefType(id, reflect.TypeOf(m_msg))
}

func gogo(t interface{}) {
	test := t.(proto.Message)
	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	var ref_type reflect.Type = dis.GetRefType(1)

	if ref_type == nil {
		log.Println("ref_type is nil")
		return
	}

	// 进行解码
	m_msg := reflect.New(ref_type.Elem()).Interface()
	proto.Unmarshal(data, m_msg.(proto.Message))

	newTest := m_msg.(*msg.Hello)

	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	log.Printf("id:%d;str:%s;", newTest.Id, newTest.Name)
}

func main() {
	dis = network.NewMsgDispatcher()
	RegisterMessage(1, &msg.Hello{})
	//	loop := network.NewEventLoop()

	//	server := network.NewServer(loop)

	//	server.Init("tcp", "127.0.0.1:80")

	//	server.GetDispatcher()

	//	loop.Loop()
	test := &msg.Hello{
		Name: "li",
		Id:   6,
	}

	gogo(test)
}
