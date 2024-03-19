package packets

import (
	"github.com/qdmc/mqtt_packet/enmu"
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
func NewPacketWithFixedHeader(f *FixedHeader) (ControlPacketInterface, error) {
	if f == nil {
		return nil, enmu.FixedEmpty
	}
	switch f.MessageType {
	case enmu.CONNECT:
		return NewConnect(f), nil
	case enmu.CONNACK:
		return NewConnAck(f), nil
	case enmu.PUBLISH:
		return NewPublish(f), nil
	case enmu.PUBACK:
		return NewPubAck(f), nil
	case enmu.PUBREC:
		return NewPubRec(f), nil
	case enmu.PUBREL:
		return NewPubRel(f), nil
	case enmu.PUBCOMP:
		return NewPubComp(f), nil
	case enmu.SUBSCRIBE:
		return NewSubscribe(f), nil
	case enmu.SUBACK:
		return NewSubAck(f), nil
	case enmu.UNSUBSCRIBE:
		return NewUnSubscribe(f), nil
	case enmu.UNSUBACK:
		return NewUnSubAck(f), nil
	case enmu.PINGREQ:
		return NewPingReq(f), nil
	case enmu.PINGRESP:
		return NewPingResp(f), nil
	case enmu.DISCONNECT:
		return NewDisconnect(f), nil
	}
	return nil, enmu.TypeError
}

func NewControlPacket(t enmu.MessageType) ControlPacketInterface {
	if t < 1 || t > 14 {
		t = enmu.PINGREQ
	}
	f := NewFixedHeader(t)
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
	default:
		return NewPingReq(f)
	}
}
