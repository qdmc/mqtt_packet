package packets

import (
	"bytes"
	"errors"
	"git.rundle.cn/bingo_queues/mqtt_packet/enmu"
	"io"
)

type TopicFilter struct {
	Topic string
	Qos   byte
}

func (t *TopicFilter) Read(r io.Reader) error {
	var err error
	t.Topic, err = decodeString(r)
	if err != nil {
		return err
	}
	t.Qos, err = decodeByte(r)
	if err != nil {
		return err
	}
	return nil
}
func (t *TopicFilter) ToBytes() []byte {
	bs := encodeString(t.Topic)
	return append(bs, t.Qos)
}

type SubscribePacket struct {
	head      *FixedHeader
	MessageID uint16
	List      []*TopicFilter
}

func NewSubscribe(f *FixedHeader) *SubscribePacket {
	return &SubscribePacket{
		head:      f,
		MessageID: 0,
		List:      []*TopicFilter{},
	}
}

func (c *SubscribePacket) MessageType() enmu.MessageType {
	return enmu.SUBSCRIBE
}

func (c *SubscribePacket) GetMessageId() uint16 {
	return c.MessageID
}

func (c *SubscribePacket) SetMessageId(id uint16) {
	c.MessageID = id
}

func (c *SubscribePacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}

func (c *SubscribePacket) Qos() byte {
	return c.GetFixedHead().Qos
}

func (c *SubscribePacket) Write(w io.Writer) (int64, error) {
	if c.List == nil || len(c.List) < 1 {
		return 0, errors.New("topics is empty")
	}
	var body bytes.Buffer
	var err error
	body.Write(encodeUint16(c.MessageID))
	for _, tf := range c.List {
		body.Write(tf.ToBytes())
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

func (c *SubscribePacket) Unpack(b io.Reader) error {
	var err error
	c.MessageID, err = decodeUint16(b)
	if err != nil {
		return err
	}
	if c.List == nil {
		c.List = []*TopicFilter{}
	}
	length := c.GetFixedHead().RemainingLength - 2
	for length > 0 {
		tf := new(TopicFilter)
		err = tf.Read(b)
		if err != nil {
			return err
		}
		length -= 2 + len(tf.Topic) + 1
	}
	return nil
}