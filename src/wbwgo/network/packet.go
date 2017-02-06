package network

type Packet struct {
	msg_id uint16

	msg_len uint16

	data []byte
}
