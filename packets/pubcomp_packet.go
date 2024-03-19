package packets

import (
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type PubCompPacket struct {
	head      *FixedHeader
	MessageID uint16
}

func NewPubComp(f *FixedHeader) *PubCompPacket {
	return &PubCompPacket{
		head:      f,
		MessageID: 0,
	}
}

func (c *PubCompPacket) MessageType() enmu.MessageType {
	return enmu.PUBCOMP
}

func (c *PubCompPacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *PubCompPacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *PubCompPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *PubCompPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *PubCompPacket) Write(w io.Writer) (int64, error) {
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

func (c *PubCompPacket) Unpack(b io.Reader) error {
	var err error
	c.MessageID, err = decodeUint16(b)
	return err
}
