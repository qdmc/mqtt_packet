package packets

import (
	"bytes"
	"encoding/binary"
	"github.com/qdmc/websocket_packet/enmu"
	"io"
)

// FixedHeader  固定报头
type FixedHeader struct {
	MessageType     enmu.MessageType // 第一字节7至4位
	Dup             bool             // 第一字节3位
	Qos             byte             // 第一字节1至2位
	Retain          bool             // 第一字节0位
	RemainingLength int              // 第2至5 字节, 剩余长度是指可变包头长度加上载荷的长度
}

// ToBuffer 写入缓冲
func (f *FixedHeader) ToBuffer() (*bytes.Buffer, error) {
	bs, err := f.ToBytes()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.Write(bs)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// ToBytes 固定报头 打码
func (f *FixedHeader) ToBytes() ([]byte, error) {
	var bs []byte
	var fistByte = 0
	fistByte = int(f.MessageType) << 4

	// 处理 标志 Flags
	if f.MessageType == enmu.PUBLISH {
		if f.Dup {
			fistByte += 1 << 3
		}
		if f.Qos == 1 {
			fistByte += 1 << 1
		} else if f.Qos == 2 {
			fistByte += 1 << 2
		}
		if f.Retain {
			fistByte += 1
		}
	} else if f.MessageType == enmu.SUBSCRIBE || f.MessageType == enmu.PUBREL || f.MessageType == enmu.UNSUBSCRIBE {
		if f.Qos == 1 {
			fistByte += 1 << 1
		}
	}
	bs = append(bs, byte(fistByte))
	if f.RemainingLength == 0 {
		return append(bs, 0), nil
	}
	lengths, err := encodeVariableByte(f.RemainingLength)
	if err != nil {
		return nil, err
	}
	bs = append(bs, lengths...)
	return bs, nil
}

// ReadFixedHeader 固定报头 解码
func ReadFixedHeader(r io.Reader) (*FixedHeader, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	b := buf[0]
	length, err := decodeVariableByte(r)
	if err != nil {
		return nil, err
	}
	h := &FixedHeader{
		MessageType:     enmu.MessageType(b >> 4),
		Dup:             (b>>3)&0x01 > 0,
		Qos:             (b >> 1) & 0x03,
		Retain:          b&0x01 > 0,
		RemainingLength: length,
	}
	err = checkQos(h.Qos)
	if err != nil {
		return nil, err
	}
	return h, nil

}

// NewFixedHeader    创建一个新的 固定报头
func NewFixedHeader(messageType enmu.MessageType) *FixedHeader {
	f := &FixedHeader{
		Qos:             0,
		Dup:             false,
		Retain:          false,
		MessageType:     messageType,
		RemainingLength: 0,
	}
	return f
}
func (f *FixedHeader) SetLength(l int) *FixedHeader {
	if l >= 0 {
		f.RemainingLength = l
	}
	return f
}
func (f *FixedHeader) SetMessageType(t enmu.MessageType) *FixedHeader {
	f.MessageType = t
	return f
}
func (f *FixedHeader) SetFixedHeaderQos(b byte) *FixedHeader {
	if f.MessageType == enmu.PUBLISH && (b == 1 || b == 0 || b == 2) {
		f.Qos = b
	}
	return f
}

func (f *FixedHeader) SetFixedHeaderRetain(b bool) *FixedHeader {
	if f.MessageType == enmu.PUBLISH || f.MessageType == enmu.PUBREL || f.MessageType == enmu.SUBSCRIBE || f.MessageType == enmu.UNSUBSCRIBE {
		f.Retain = b
	}
	return f
}

func (f *FixedHeader) SetFixedHeaderDup(b bool) *FixedHeader {
	if f.MessageType == enmu.PUBLISH {
		f.Dup = b
	}
	return f
}
func decodeBytes(b io.Reader) ([]byte, error) {
	fieldLength, err := decodeUint16(b)
	if err != nil {
		return nil, err
	}

	field := make([]byte, fieldLength)
	_, err = b.Read(field)
	if err != nil {
		return nil, err
	}

	return field, nil
}
func decodeString(r io.Reader) (string, error) {
	return readString(r)
}

func decodeByte(b io.Reader) (byte, error) {
	return readByte(b)
}
func decodeUint16(b io.Reader) (uint16, error) {
	return readUin16(b)
}
func encodeUint16(num uint16) []byte {
	bytesResult := make([]byte, 2)
	binary.BigEndian.PutUint16(bytesResult, num)
	return bytesResult
}

func encodeString(field string) []byte {
	return encodeBytes([]byte(field))
}
func encodeBytes(field []byte) []byte {
	fieldLength := make([]byte, 2)
	binary.BigEndian.PutUint16(fieldLength, uint16(len(field)))
	return append(fieldLength, field...)
}

func boolToByte(b bool) byte {
	switch b {
	case true:
		return 1
	default:
		return 0
	}
}

func checkQos(qos byte) error {
	if qos > 2 {
		return enmu.QosError
	}
	return nil
}

func checkTopicName(name string) error {
	if name == "" {
		return enmu.TopicError
	}
	return nil
}

func checkSubAckPayload(payload []byte) error {
	if payload == nil || len(payload) < 1 {
		return enmu.SubAckPayloadEmpty
	}
	for i, _ := range payload {
		b := payload[i]
		if b == 0x00 || b == 0x01 || b == 0x02 || b == 0x80 {
			continue
		} else {
			return enmu.SubAckPayloadError
		}
	}
	return nil
}
