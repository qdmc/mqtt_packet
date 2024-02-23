package packets

import (
	"git.rundle.cn/bingo_queues/mqtt_packet/enmu"
	"io"
)

/*
ControlPacketInterface        报文通用接口
  - MessageType()              返回报文类型
  - GetMessageId()             返回报文Id
  - SetMessageId(id uint16)    设置消息Id
  - GetFixedHead()             返回报文的固定报头
  - Qos()                      返回报文的Qos
  - Write(w io.Writer)         报文写入
  - Unpack(b io.Reader)        报文读取
*/
type ControlPacketInterface interface {
	MessageType() enmu.MessageType
	GetMessageId() uint16
	SetMessageId(id uint16)
	GetFixedHead() *FixedHeader
	Qos() byte
	Write(w io.Writer) (int64, error)
	Unpack(b io.Reader) error
}

// NewPacketWithFixedHeader   根据固定报头生成报文
func NewPacketWithFixedHeader(f *FixedHeader) ControlPacketInterface {
	if f == nil {
		return nil
	}
	switch f.MessageType {
	case enmu.CONNECT:
		return NewConnect(f)
	case enmu.CONNACK:
		return NewConnAck(f)
	case enmu.PUBLISH:
		return NewPublish(f)
	case enmu.PUBACK:
		return NewPubAck(f)
	case enmu.PUBREC:
		return NewPubRec(f)
	case enmu.PUBREL:
		return NewPubRel(f)
	case enmu.PUBCOMP:
		return NewPubComp(f)
	case enmu.SUBSCRIBE:
		return NewSubscribe(f)
	case enmu.SUBACK:
		return NewSubAck(f)
	case enmu.UNSUBSCRIBE:
		return NewUnSubscribe(f)
	case enmu.UNSUBACK:
		return NewUnSubAck(f)
	case enmu.PINGREQ:
		return NewPingReq(f)
	case enmu.PINGRESP:
		return NewPingResp(f)
	case enmu.DISCONNECT:
		return NewDisconnect(f)
	}
	return nil
}
