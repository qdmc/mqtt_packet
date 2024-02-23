package packets

import (
	"git.rundle.cn/bingo_queues/mqtt_packet/enmu"
	"io"
)

type PubAckPacket struct {
	head      *FixedHeader
	MessageID uint16
}

func NewPubAck(f *FixedHeader) *PubAckPacket {
	return &PubAckPacket{
		head:      f,
		MessageID: 0,
	}
}

func (c *PubAckPacket) MessageType() enmu.MessageType {
	return enmu.PUBACK
}

func (c *PubAckPacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *PubAckPacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *PubAckPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *PubAckPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *PubAckPacket) Write(w io.Writer) (int64, error) {
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

func (c *PubAckPacket) Unpack(b io.Reader) error {
	var err error
	c.MessageID, err = decodeUint16(b)
	return err
}
