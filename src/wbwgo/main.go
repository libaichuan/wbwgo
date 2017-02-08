package main

import (
	"log"
	"wbwgo/msg"
	"wbwgo/network"
)

var ch chan int

func main() {
	ch = make(chan int)

	loop := network.NewEventLoop()

	server := network.NewServer(loop)

	server.Init("tcp", "127.0.0.1:8000")

	network.RegisterMessage(server.GetDispatcher(), 1, &msg.Hello{}, func(ses *network.Session, f interface{}) {
		cur_msg := f.(*msg.Hello)
		log.Println("recv new msg id:%d,name:%s\n", cur_msg.Id, cur_msg.Name)
		ch <- 1
	})

	loop.Loop()

	<-ch
}
