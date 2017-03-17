package main

import (
	//"fmt"
	"time"
	"wbwgo/game"
)

func main() {

	o := new(game.MapObjectBase)

	o.ComputeSpeed(3.0, 4.0, 10.0)

	timer := time.NewTicker(100 * time.Millisecond)

	for {
		_, ok := <-timer.C
		if ok {
			o.Move()
			o.Print()
		}
	}
}
