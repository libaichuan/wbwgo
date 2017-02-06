package network

//事件循环

/*
eventloop负责 IO多线程收到的data通过eventloop传给单线程的业务
			  timer内部调用

那么这里就要设计咋区分这两种调用  要是外部的data过来，通过Packet的数据反序列化
 							timer触发的话直接调.

加入的时候外部Proto事件告知dispatcher,以及session和packet，内部事件dispatcher为Nil，data传func直接调用。
*/

type LoopData struct {
	dispatcher MsgDispatcher

	data interface{}
}

type EventLoop struct {
	queue chan LoopData
}

func (self *EventLoop) Loop() {
	for q := range self.queue {

		if q.dispatcher != nil {
			d := q.data()
			if f, ok := d.(func()); ok {
				f()
			}
		} else {
			q.dispatcher.OnMessage(q.data)
		}

	}
}

func (self *EventLoop) AddInLoop(dp MsgDispatcher, data interface{}) {
	self.queue <- LoopData{dispatcher: dp, data: data}
}

func NewEventLoop() *EventLoop {
	self := &EventLoop{
		queue: make(chan LoopData, 5),
	}

	return self
}
