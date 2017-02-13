package network

import (
	//"log"
	"net"
)

type Connect struct {
	conn        net.Conn
	event_loop  *EventLoop
	dispatcher  *MsgDispatcher
	ses_manager *SessionManager
	closeSignal chan bool
}

func (self *Connect) GetDispatcher() *MsgDispatcher {
	return self.dispatcher
}

func (self *Connect) Start(addr string) {
	go self.DoConnet(addr)
}
func (self *Connect) DoConnet(addr string) {
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return
		}

		//log.Println("client connet.......")

		self.conn = conn

		ses := NewSession(conn, self.dispatcher, self.event_loop)

		ses.OnClose = func() {
			self.ses_manager.RemoveSessionById(ses.id)
			self.closeSignal <- true
		}

		self.ses_manager.AddSession(ses)

		self.event_loop.AddInLoop(self.dispatcher, NewEventData(ses, &Packet{
			msg_id: 2,
		}))

		<-self.closeSignal
	}
}

func NewClient(event_loop *EventLoop) *Connect {
	self := new(Connect)
	self.event_loop = event_loop
	self.dispatcher = NewMsgDispatcher()
	self.ses_manager = NewSessionManager()

	return self
}
