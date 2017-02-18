package benchmark

import (
	"testing"
	"time"
	"wbwgo/benchmark"
	//"wbwgo/common"
	"log"
	"os"
	"runtime/pprof"
	"wbwgo/msg"
	"wbwgo/network"
)

var ch chan int

func NewServer() {
	loop := network.NewEventLoop()

	qps, _ := benchmark.NewQpsMeter(loop)

	server := network.NewServer(loop)

	server.Init("tcp", "127.0.0.1:8000")

	network.RegisterMessage(server.GetDispatcher(), 1, &msg.Hello{}, func(ses *network.Session, f interface{}) {
		qps.AddQps()

		cur_msg := f.(*msg.Hello)
		//log.Println("sever recv new msg id:%d,name:%s\n", cur_msg.Id, cur_msg.Name)
		ses.Send(cur_msg)
	})

	loop.Loop()

	time.AfterFunc(time.Second*10, func() {
		ch <- 1
	})
}

func NewClient() {
	loop := network.NewEventLoop()

	client := network.NewClient(loop)

	client.Start("127.0.0.1:8000")

	network.RegisterMessage(client.GetDispatcher(), 1, &msg.Hello{}, func(ses *network.Session, f interface{}) {
		cur_msg := f.(*msg.Hello)
		//log.Println("client recv new msg id:%d,name:%s\n", cur_msg.Id, cur_msg.Name)
		ses.Send(cur_msg)
	})

	network.RegisterMessage(client.GetDispatcher(), 2, &msg.OnSessionConnet{}, func(ses *network.Session, f interface{}) {
		//log.Println("client recv new msg OnSessionConnet")
		msg_s := &msg.Hello{}
		msg_s.Id = 1
		msg_s.Name = "gogo"
		ses.Send(msg_s)
	})

	loop.Loop()
}

func GoGo_IO() {
	client_num := 100

	ch = make(chan int)

	//common.ConsoleStart()

	NewServer()

	for i := 0; i < client_num; i++ {
		go NewClient()
	}

	<-ch
}

func BenchmarkIO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GoGo_IO()
	}
}
