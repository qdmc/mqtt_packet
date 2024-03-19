package packets

import (
	"bytes"
	"github.com/qdmc/websocket_packet/enmu"
	"io"
)

type ConnectPacket struct {
	head             *FixedHeader
	ProtocolName     string
	ProtocolVersion  byte
	CleanSession     bool
	WillFlag         bool
	WillQos          byte
	WillRetain       bool
	UsernameFlag     bool
	PasswordFlag     bool
	ReservedBit      byte
	Keepalive        uint16
	ClientIdentifier string
	WillTopic        string
	WillMessage      []byte
	Username         string
	Password         []byte
}

func NewConnect(f *FixedHeader) *ConnectPacket {
	return &ConnectPacket{
		head:             f,
		ProtocolName:     "",
		ProtocolVersion:  0,
		CleanSession:     false,
		WillFlag:         false,
		WillQos:          0,
		WillRetain:       false,
		UsernameFlag:     false,
		PasswordFlag:     false,
		ReservedBit:      0,
		Keepalive:        0,
		ClientIdentifier: "",
		WillTopic:        "",
		WillMessage:      nil,
		Username:         "",
		Password:         nil,
	}
}
func (c *ConnectPacket) SetMessageId(id uint16) {

}
func (c *ConnectPacket) GetFixedHead() *FixedHeader {
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
	}
	return c.head
}
func (c *ConnectPacket) Qos() byte {
	return c.GetFixedHead().Qos
}
func (c *ConnectPacket) MessageType() enmu.MessageType {
	return enmu.CONNECT
}

func (c *ConnectPacket) GetMessageId() uint16 {
	return 0
}

func (c *ConnectPacket) Write(w io.Writer) (int64, error) {
	var body bytes.Buffer
	var err error
	body.Write(encodeString(c.ProtocolName))
	if c.ProtocolVersion != 3 && c.ProtocolVersion != 4 {
		c.ProtocolVersion = 3
	}
	body.WriteByte(c.ProtocolVersion)
	body.WriteByte(boolToByte(c.CleanSession)<<1 | boolToByte(c.WillFlag)<<2 | c.WillQos<<3 | boolToByte(c.WillRetain)<<5 | boolToByte(c.PasswordFlag)<<6 | boolToByte(c.UsernameFlag)<<7)
	body.Write(encodeUint16(c.Keepalive))
	body.Write(encodeString(c.ClientIdentifier))
	if c.WillFlag {
		body.Write(encodeString(c.WillTopic))
		body.Write(encodeBytes(c.WillMessage))
	}
	if c.UsernameFlag {
		body.Write(encodeString(c.Username))
	}
	if c.PasswordFlag {
		body.Write(encodeBytes(c.Password))
	}
	if c.head == nil {
		c.head = NewFixedHeader(c.MessageType())
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

func (c *ConnectPacket) Unpack(b io.Reader) error {
	var err error
	c.ProtocolName, err = decodeString(b)
	if err != nil {
		return err
	}
	c.ProtocolVersion, err = decodeByte(b)
	if err != nil {
		return err
	}
	options, err := decodeByte(b)
	if err != nil {
		return err
	}
	c.ReservedBit = 1 & options
	c.CleanSession = 1&(options>>1) > 0
	c.WillFlag = 1&(options>>2) > 0
	c.WillQos = 3 & (options >> 3)
	c.WillRetain = 1&(options>>5) > 0
	c.PasswordFlag = 1&(options>>6) > 0
	c.UsernameFlag = 1&(options>>7) > 0
	c.Keepalive, err = decodeUint16(b)
	if err != nil {
		return err
	}
	c.ClientIdentifier, err = decodeString(b)
	if err != nil {
		return err
	}
	if c.WillFlag {
		c.WillTopic, err = decodeString(b)
		if err != nil {
			return err
		}
		c.WillMessage, err = decodeBytes(b)
		if err != nil {
			return err
		}
	}
	if c.UsernameFlag {
		c.Username, err = decodeString(b)
		if err != nil {
			return err
		}
	}
	if c.PasswordFlag {
		c.Password, err = decodeBytes(b)
		if err != nil {
			return err
		}
	}

	return nil
}
