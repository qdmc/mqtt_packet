package mqtt_packet

import (
	"bytes"
	"git.rundle.cn/bingo_queues/mqtt_packet/enmu"
	"git.rundle.cn/bingo_queues/mqtt_packet/packets"
	"io"
)

// MessageType   消息类型
type MessageType = enmu.MessageType

// FixedHeader    固定报头
type FixedHeader = packets.FixedHeader

func NewFixedHeader(t MessageType) *FixedHeader {
	return packets.NewFixedHeader(t)
}

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
type ControlPacketInterface = packets.ControlPacketInterface

// NewControlPacketWithFixedHead           根据固定报头生成报文
func NewControlPacketWithFixedHead(f *FixedHeader) (ControlPacketInterface, error) {
	return packets.NewPacketWithFixedHeader(f)
}

// NewControlPacket                        生成报文,默认为PingReq报文
func NewControlPacket(t enmu.MessageType) ControlPacketInterface {
	return packets.NewControlPacket(t)
}

// NewFixedHead               创建一个新的 固定报头
func NewFixedHead(t MessageType) *FixedHeader {
	return packets.NewFixedHeader(t)
}

// ReadOnce      读取一个报文
func ReadOnce(r io.Reader) (ControlPacketInterface, error) {
	h, err := packets.ReadFixedHeader(r)
	if err != nil {
		return nil, err
	}
	cp, err := packets.NewPacketWithFixedHeader(h)
	if err != nil {
		return nil, err
	}
	err = cp.Unpack(r)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

// ReadStream            读取字节流
func ReadStream(bs []byte) (list []ControlPacketInterface, lastBytes []byte, err error) {
	var header *FixedHeader
	var cp ControlPacketInterface
	for {
		if len(bs) == 0 || len(bs) == 1 {
			return list, bs, nil
		}
		headBs := bs[:2]
		bs = bs[2:]
		header, err = packets.ReadFixedHeader(bytes.NewBuffer(headBs))
		if err != nil {
			return list, bs, err
		}
		cp, err = NewControlPacketWithFixedHead(header)
		if err != nil {
			return list, bs, err
		}
		//println("---- RemainingLength:　", header.RemainingLength, "　　　bsLen: ", len(bs), "  messageType: ", header.MessageType)
		if header.RemainingLength == 0 {
			list = append(list, cp)
			continue
		} else {
			if header.RemainingLength <= len(bs) {
				packetBs := bs[:header.RemainingLength]
				bs = bs[header.RemainingLength:]
				err = cp.Unpack(bytes.NewBuffer(packetBs))
				if err != nil {
					return list, bs, err
				}
				list = append(list, cp)
				continue
			} else {
				return list, bs, nil
			}
		}
	}
}

type ConnAckPacket = packets.ConnAckPacket

func NewConnAck(f *FixedHeader) *ConnAckPacket {
	return packets.NewConnAck(f)
}

type ConnectPacket = packets.ConnectPacket

func NewConnect(f *FixedHeader) *ConnectPacket {
	return packets.NewConnect(f)
}

type DisconnectPacket = packets.DisconnectPacket

func NewDisconnect(f *FixedHeader) *DisconnectPacket {
	return packets.NewDisconnect(f)
}

type PingReqPacket = packets.PingReqPacket

func NewPingReq(f *FixedHeader) *PingReqPacket {
	return packets.NewPingReq(f)
}

type PingRespPacket = packets.PingRespPacket

func NewPingResp(f *FixedHeader) *PingRespPacket {
	return packets.NewPingResp(f)
}

type PubAckPacket = packets.PubAckPacket

func NewPubAck(f *FixedHeader) *PubAckPacket {
	return packets.NewPubAck(f)
}

type PubCompPacket = packets.PubCompPacket

func NewPubComp(f *FixedHeader) *PubCompPacket {
	return packets.NewPubComp(f)
}

type PublishPacket = packets.PublishPacket

func NewPublish(f *FixedHeader) *PublishPacket {
	return packets.NewPublish(f)
}

type PubRecPacket = packets.PubRecPacket

func NewPubRec(f *FixedHeader) *PubRecPacket {
	return packets.NewPubRec(f)
}

type PubRelPacket = packets.PubRelPacket

func NewPubRel(f *FixedHeader) *PubRelPacket {
	return packets.NewPubRel(f)
}

type SubAckPacket = packets.SubAckPacket

func NewSubAck(f *FixedHeader) *SubAckPacket {
	return packets.NewSubAck(f)
}

type SubscribePacket = packets.SubscribePacket

func NewSubscribe(f *FixedHeader) *SubscribePacket {
	return packets.NewSubscribe(f)
}

type UnSubAckPacket = packets.UnSubAckPacket

func NewUnSubAck(f *FixedHeader) *UnSubAckPacket {
	return packets.NewUnSubAck(f)
}

type UnSubscribePacket = packets.UnSubscribePacket

func NewUnSubscribe(f *FixedHeader) *UnSubscribePacket {
	return packets.NewUnSubscribe(f)
}
