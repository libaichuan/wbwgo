package misc

import (
	//"log"
	"time"
	"wbwgo/network"
)

type Timer struct {
	ticker     *time.Ticker
	close_sign chan bool
}

func (self *Timer) Close() {
	self.close_sign <- true
}

func NewTimer(event_loop *network.EventLoop, dur time.Duration, f func()) *Timer {
	timer := &Timer{
		ticker:     time.NewTicker(dur),
		close_sign: make(chan bool),
	}

	go func() {
		for {
			select {
			case <-timer.ticker.C:
				event_loop.AddInLoop(nil, f)
			case <-timer.close_sign:
				timer.ticker.Stop()
				break
			}
		}
	}()

	return timer
}
