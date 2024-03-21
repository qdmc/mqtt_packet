package packets

import (
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type PingRespPacket struct {
	head *FixedHeader
}

func NewPingResp(f *FixedHeader) *PingRespPacket {
	return &PingRespPacket{
		head: f,
	}
}

func (c *PingRespPacket) MessageType() enmu.MessageType {
	return enmu.PINGRESP
}

func (c *PingRespPacket) GetMessageId() uint16 {
	return 0
}

func (c *PingRespPacket) SetMessageId(id uint16) {
}

func (c *PingRespPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *PingRespPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *PingRespPacket) Write(w io.Writer) (int64, error) {
	head := c.GetFixedHead()
	buf, err := head.ToBuffer()
	if err != nil {
		return 0, err
	}
	return buf.WriteTo(w)
}

func (c *PingRespPacket) Unpack(b io.Reader) (int, error) {
	return 0, nil
}
