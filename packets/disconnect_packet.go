package packets

import (
	"github.com/qdmc/websocket_packet/enmu"
	"io"
)

type DisconnectPacket struct {
	head *FixedHeader
}

func NewDisconnect(f *FixedHeader) *DisconnectPacket {
	return &DisconnectPacket{head: f}
}

func (c *DisconnectPacket) MessageType() enmu.MessageType {
	return enmu.DISCONNECT
}

func (c *DisconnectPacket) GetMessageId() uint16 {
	return 0
}

func (c *DisconnectPacket) SetMessageId(id uint16) {
}

func (c *DisconnectPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *DisconnectPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *DisconnectPacket) Write(w io.Writer) (int64, error) {
	head := c.GetFixedHead()
	buf, err := head.ToBuffer()
	if err != nil {
		return 0, err
	}
	return buf.WriteTo(w)
}

func (c *DisconnectPacket) Unpack(b io.Reader) error {
	return nil
}
