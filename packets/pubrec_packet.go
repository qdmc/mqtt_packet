package packets

import (
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type PubRecPacket struct {
	head      *FixedHeader
	MessageID uint16
}

func NewPubRec(f *FixedHeader) *PubRecPacket {
	return &PubRecPacket{
		head:      f,
		MessageID: 0,
	}
}

func (c *PubRecPacket) MessageType() enmu.MessageType {
	return enmu.PUBREC
}

func (c *PubRecPacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *PubRecPacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *PubRecPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *PubRecPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *PubRecPacket) Write(w io.Writer) (int64, error) {
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

func (c *PubRecPacket) Unpack(b io.Reader) (int, error) {
	var err error
	c.MessageID, err = decodeUint16(b)
	return 2, err
}
