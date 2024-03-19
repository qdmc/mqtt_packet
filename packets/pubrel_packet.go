package packets

import (
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type PubRelPacket struct {
	head      *FixedHeader
	MessageID uint16
}

func NewPubRel(f *FixedHeader) *PubRelPacket {
	return &PubRelPacket{
		head:      f,
		MessageID: 0,
	}
}

func (c *PubRelPacket) MessageType() enmu.MessageType {
	return enmu.PUBREL
}

func (c *PubRelPacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *PubRelPacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *PubRelPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *PubRelPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *PubRelPacket) Write(w io.Writer) (int64, error) {
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

func (c *PubRelPacket) Unpack(b io.Reader) error {
	var err error
	c.MessageID, err = decodeUint16(b)
	return err
}
