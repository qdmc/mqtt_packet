package packets

import (
	"bytes"
	"github.com/qdmc/websocket_packet/enmu"
	"io"
)

type ConnAckPacket struct {
	head           *FixedHeader
	SessionPresent bool
	ReturnCode     byte
}

func NewConnAck(f *FixedHeader) *ConnAckPacket {
	return &ConnAckPacket{
		head:           f,
		SessionPresent: false,
		ReturnCode:     0,
	}
}

func (c *ConnAckPacket) MessageType() enmu.MessageType {
	return enmu.CONNACK
}

func (c *ConnAckPacket) GetMessageId() uint16 {
	return 0
}

func (c *ConnAckPacket) SetMessageId(id uint16) {

}
func (c *ConnAckPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}
func (c *ConnAckPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *ConnAckPacket) Write(w io.Writer) (int64, error) {
	var body bytes.Buffer
	var err error
	body.WriteByte(boolToByte(c.SessionPresent))
	body.WriteByte(c.ReturnCode)
	head := c.GetFixedHead()
	head.MessageType = c.MessageType()
	head.RemainingLength = 2
	buf, err := head.ToBuffer()
	if err != nil {
		return 0, err
	}
	_, err = buf.Write(body.Bytes())
	if err != nil {
		return 0, err
	}
	return buf.WriteTo(w)
}

func (c *ConnAckPacket) Unpack(b io.Reader) error {
	flags, err := decodeByte(b)
	if err != nil {
		return err
	}
	c.SessionPresent = 1&flags > 0
	c.ReturnCode, err = decodeByte(b)
	return err
}
