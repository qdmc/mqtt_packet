package packets

import (
	"bytes"
	"fmt"
	"github.com/qdmc/websocket_packet/enmu"
	"io"
)

type PublishPacket struct {
	head      *FixedHeader
	TopicName string
	MessageID uint16
	Payload   []byte
}

func NewPublish(f *FixedHeader) *PublishPacket {
	return &PublishPacket{
		head:      f,
		TopicName: "",
		MessageID: 0,
		Payload:   nil,
	}
}

func (c *PublishPacket) MessageType() enmu.MessageType {
	return enmu.PUBLISH
}

func (c *PublishPacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *PublishPacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *PublishPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *PublishPacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *PublishPacket) Write(w io.Writer) (int64, error) {
	var body bytes.Buffer
	var err error
	body.Write(encodeString(c.TopicName))
	if c.Qos() > 0 {
		body.Write(encodeUint16(c.MessageID))
	}
	body.Write(c.Payload)

	head := c.GetFixedHead()
	head.MessageType = c.MessageType()
	head.RemainingLength = body.Len()
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

func (c *PublishPacket) Unpack(b io.Reader) error {
	var payloadLength = c.GetFixedHead().RemainingLength
	var err error
	c.TopicName, err = decodeString(b)
	if err != nil {
		return err
	}
	if c.Qos() > 0 {
		c.MessageID, err = decodeUint16(b)
		if err != nil {
			return err
		}
		payloadLength -= len(c.TopicName) + 4
	} else {
		payloadLength -= len(c.TopicName) + 2
	}
	if payloadLength < 0 {
		return fmt.Errorf("error unpacking publish, payload length < 0")
	}
	c.Payload = make([]byte, payloadLength)
	_, err = b.Read(c.Payload)
	return err
}
