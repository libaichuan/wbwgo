package network

type EventData struct {
	ses *Session

	p Packet
}

func NewEventData(ses *Session, p Packet) *EventData {
	self := &EventData{
		ses: ses,
		p:   p,
	}

	return self
}
