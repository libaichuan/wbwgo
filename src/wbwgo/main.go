package main

import (
	"log"
	"time"
	"wbwgo/network"
)

var s_manager *network.SessionManager

func DoSome() {
	for i := 0; i < 10000; i++ {
		go func() {
			s := network.NewSession(nil)
			s_manager.AddSession(s)
		}()
	}
}

func main() {
	s_manager = network.NewSessionManager()

	DoSome()

	time.Sleep(10 * time.Second)

	log.Printf("total num %d\n", s_manager.GetSessionCount())

	time.Sleep(time.Second)

}
