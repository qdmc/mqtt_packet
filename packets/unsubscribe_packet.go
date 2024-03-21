package packets

import (
	"bytes"
	"errors"
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type UnSubscribePacket struct {
	head      *FixedHeader
	MessageID uint16
	Topics    []string
}

func NewUnSubscribe(f *FixedHeader) *UnSubscribePacket {
	return &UnSubscribePacket{
		head:      f,
		MessageID: 0,
		Topics:    nil,
	}
}

func (c *UnSubscribePacket) MessageType() enmu.MessageType {
	return enmu.UNSUBSCRIBE
}

func (c *UnSubscribePacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *UnSubscribePacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *UnSubscribePacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *UnSubscribePacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *UnSubscribePacket) Write(w io.Writer) (int64, error) {
	if c.Topics == nil || len(c.Topics) < 1 {
		return 0, errors.New("UnSubscribe.Topics is empty")
	}
	var body bytes.Buffer
	var err error
	_, err = body.Write(encodeUint16(c.MessageID))
	if err != nil {
		return 0, err
	}
	for i, _ := range c.Topics {
		title := c.Topics[i]
		err = checkTopicName(title)
		if err != nil {
			return 0, err
		}
		_, err = body.Write(encodeString(title))
		if err != nil {
			return 0, err
		}
	}
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

func (c *UnSubscribePacket) Unpack(b io.Reader) (int, error) {
	var err error
	var readLen, rl int
	if c.head == nil {
		return readLen, enmu.FixedEmpty
	}
	c.MessageID, err = decodeUint16(b)
	if err != nil {
		return readLen, err
	}
	readLen += 2
	length := c.head.RemainingLength - 2
	c.Topics = []string{}
	var topic string
	for length > 0 {
		rl, topic, err = decodeString(b)
		if err != nil {
			return readLen, err
		}
		err = checkTopicName(topic)
		if err != nil {
			return readLen, err
		}
		readLen += rl
		c.Topics = append(c.Topics, topic)
		length -= 2 + len([]byte(topic))
	}
	return readLen, err
}
