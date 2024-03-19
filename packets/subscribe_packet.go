package packets

import (
	"bytes"
	"errors"
	"github.com/qdmc/mqtt_packet/enmu"
	"io"
)

type TopicFilter struct {
	Topic string
	Qos   byte
}

func (t *TopicFilter) check() error {
	err := checkTopicName(t.Topic)
	if err != nil {
		return err
	}
	err = checkQos(t.Qos)
	if err != nil {
		return err
	}
	return nil
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
	var topicsBs []byte
	for _, tf := range c.List {
		if tf == nil {
			continue
		}
		err = tf.check()
		if err != nil {
			return 0, err
		}
		//body.Write(tf.ToBytes())
		topicsBs = append(topicsBs, tf.ToBytes()...)
	}
	if len(topicsBs) == 0 {
		return 0, enmu.TopicError
	}
	body.Write(topicsBs)
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
	if c.head == nil {
		return enmu.FixedEmpty
	}
	length := c.head.RemainingLength - 2
	c.List = []*TopicFilter{}
	for length > 0 {
		tf := new(TopicFilter)
		err = tf.Read(b)
		if err != nil {
			return err
		}
		err = tf.check()
		if err != nil {
			return err
		}
		c.List = append(c.List, tf)
		length -= 2 + len([]byte(tf.Topic)) + 1
	}
	if len(c.List) == 0 {
		return enmu.TopicsEmpty
	}
	return nil
}
