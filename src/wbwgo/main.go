package main

import (
	"wbwgo/network"
)

func main() {

	loop := NewEventLoop()

	server := NewServer(loop)

}
