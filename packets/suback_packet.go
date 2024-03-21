package packets

import (
	"bytes"
	"github.com/qdmc/mqtt_packet/enmu"
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
	err = checkSubAckPayload(c.ReturnCodes)
	if err != nil {
		return 0, err
	}
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

func (c *SubAckPacket) Unpack(b io.Reader) (int, error) {
	var readLen, rl int
	if c.head == nil {
		return readLen, enmu.FixedEmpty
	}
	if c.head.RemainingLength < 2 {
		return readLen, enmu.RemainingLengthErr
	}
	var err error
	c.MessageID, err = decodeUint16(b)
	if err != nil {
		return readLen, err
	}
	readLen += 2
	codeBs := make([]byte, c.head.RemainingLength-2)
	rl, err = b.Read(codeBs)
	if err != nil {
		return readLen, err
	}
	readLen += rl
	err = checkSubAckPayload(codeBs)
	if err != nil {
		return readLen, err
	}
	c.ReturnCodes = codeBs
	return readLen, nil
}
