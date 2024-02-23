package packets

import (
	"bytes"
	"git.rundle.cn/bingo_queues/mqtt_packet/enmu"
	"io"
)

type SubAckPacket struct {
	head        *FixedHeader
	MessageID   uint16
	ReturnCodes []byte
}

func NewSubAck(f *FixedHeader) *SubAckPacket {
	return &SubAckPacket{
		head:        f,
		MessageID:   0,
		ReturnCodes: nil,
	}
}

func (c *SubAckPacket) MessageType() enmu.MessageType {
	return enmu.SUBACK
}

func (c *SubAckPacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *SubAckPacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *SubAckPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *SubAckPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *SubAckPacket) Write(w io.Writer) (int64, error) {
	var body bytes.Buffer
	var err error
	body.Write(encodeUint16(c.MessageID))
	body.Write(c.ReturnCodes)
	head := c.GetFixedHead()
	head.RemainingLength = body.Len()
	head.MessageType = c.MessageType()
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

func (c *SubAckPacket) Unpack(b io.Reader) error {
	if c.head == nil {
		return enmu.FixedEmpty
	}
	if c.head.RemainingLength < 2 {
		return enmu.RemainingLengthErr
	}
	var err error
	c.MessageID, err = decodeUint16(b)
	if err != nil {
		return err
	}
	codeBs := make([]byte, c.head.RemainingLength-2)
	_, err = b.Read(codeBs)
	return err
}
