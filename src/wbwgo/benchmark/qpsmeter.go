package benchmark

import (
	"log"
	"sync"
	"time"
	"wbwgo/misc"
	"wbwgo/network"
)

type QpsMeter struct {
	cur_qps   int
	total_qps int

	cur_index int
	mutex     sync.Mutex
}

func (self *QpsMeter) AddQps() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.cur_qps++
}

func (self *QpsMeter) OverIndex() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.cur_index++

	self.total_qps += self.cur_qps

	self.cur_qps = 0
}

func (self *QpsMeter) PrintInfo() int {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.cur_index == 0 {
		return 0
	}

	return self.total_qps / self.cur_index
}

func NewQpsMeter(even_loop *network.EventLoop) (*QpsMeter, *misc.Timer) {
	p := &QpsMeter{}

	timer := misc.NewTimer(even_loop, time.Second, func() {
		log.Printf("cur_qps:%d,index:%d", p.cur_qps, p.cur_index)
		p.OverIndex()
		//log.Printf("acc:%d,index:%d", p.PrintInfo(), p.cur_index)
	})

	return p, timer
}
