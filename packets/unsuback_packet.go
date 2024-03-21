package packets

import (
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type UnSubAckPacket struct {
	head      *FixedHeader
	MessageID uint16
}

func NewUnSubAck(f *FixedHeader) *UnSubAckPacket {
	return &UnSubAckPacket{
		head:      f,
		MessageID: 0,
	}
}

func (c *UnSubAckPacket) MessageType() enmu.MessageType {
	return enmu.UNSUBACK
}

func (c *UnSubAckPacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *UnSubAckPacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *UnSubAckPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *UnSubAckPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *UnSubAckPacket) Write(w io.Writer) (int64, error) {
	var err error
	head := c.GetFixedHead()
	head.MessageType = c.MessageType()
	head.RemainingLength = 2
	buf, err := head.ToBuffer()
	if err != nil {
		return 0, err
	}
	buf.Write(encodeUint16(c.MessageID))
	return buf.WriteTo(w)
}

func (c *UnSubAckPacket) Unpack(b io.Reader) (int, error) {
	var err error
	c.MessageID, err = decodeUint16(b)
	return 2, err
}
