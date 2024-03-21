package packets

import (
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type PingReqPacket struct {
	head *FixedHeader
}

func NewPingReq(f *FixedHeader) *PingReqPacket {
	return &PingReqPacket{head: f}
}

func (c *PingReqPacket) MessageType() enmu.MessageType {
	return enmu.PINGREQ
}

func (c *PingReqPacket) GetMessageId() uint16 {
	return 0
}

func (c *PingReqPacket) SetMessageId(id uint16) {
}

func (c *PingReqPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *PingReqPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *PingReqPacket) Write(w io.Writer) (int64, error) {
	head := c.GetFixedHead()
	buf, err := head.ToBuffer()
	if err != nil {
		return 0, err
	}
	return buf.WriteTo(w)
}

func (c *PingReqPacket) Unpack(b io.Reader) (int, error) {
	return 0, nil
}
