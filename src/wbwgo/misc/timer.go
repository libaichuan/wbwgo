package misc

import (
	"time"
	"wbwgo/network"
)

type Timer struct {
}

func NewTimer(event_loop *network.EventLoop, dur time.Duration, f func()) {
	ticker := time.NewTicker(dur)
	go func() {
		for {
			_, ok := <-ticker.C
			if ok {
				event_loop.AddInLoop(nil, f)
			} else {
				break
			}
		}
	}()
}
