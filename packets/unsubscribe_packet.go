package packets

import (
	"bytes"
	"errors"
	"git.rundle.cn/bingo_queues/mqtt_packet/enmu"
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
		_, err = body.Write(encodeString(c.Topics[i]))
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

func (c *UnSubscribePacket) Unpack(b io.Reader) error {
	var err error
	c.MessageID, err = decodeUint16(b)
	if err != nil {
		return err
	}

	for topic, err := decodeString(b); err == nil && topic != ""; topic, err = decodeString(b) {
		c.Topics = append(c.Topics, topic)
	}
	return err
}
