package test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"git.rundle.cn/bingo_queues/mqtt_packet"
	"git.rundle.cn/bingo_queues/mqtt_packet/enmu"
	"testing"
)

var hexStrMap map[enmu.MessageType]string

func init() {
	hexStrMap = map[enmu.MessageType]string{
		enmu.CONNECT:     "102C00044D51545404C2003C000A636C69656E7469642F31000A757365726E616D652F31000870617373776F7264",
		enmu.CONNACK:     "20020000",
		enmu.PUBLISH:     "32100005746F70696300016D657373616765",
		enmu.PUBACK:      "40020001",
		enmu.PUBREC:      "50020001",
		enmu.PUBREL:      "60020001",
		enmu.PUBCOMP:     "70020001",
		enmu.SUBSCRIBE:   "820A00010005746F70696300",
		enmu.SUBACK:      "9003000100",
		enmu.UNSUBSCRIBE: "A20900100005746F706963",
		enmu.UNSUBACK:    "B0020010",
		enmu.PINGREQ:     "C000",
		enmu.PINGRESP:    "D000",
		enmu.DISCONNECT:  "E000",
	}
}
func getPacketBytes(t enmu.MessageType) ([]byte, error) {
	if str, ok := hexStrMap[t]; ok {
		return hex.DecodeString(str)
	} else {
		return nil, enmu.TypeError
	}
}
func writePacket(cp mqtt_packet.ControlPacketInterface, packetBytesLength int) error {
	writeBuf := bytes.NewBuffer([]byte{})
	length, err := cp.Write(writeBuf)
	if err != nil {
		return errors.New(fmt.Sprintf("packetWriteErr: %s", err.Error()))
	}
	if int(length) != packetBytesLength {
		return errors.New(fmt.Sprintf("packetWriteLenErr: %d ,packetBsLne: %d", length, packetBytesLength))
	}
	if writeBuf.Len() != packetBytesLength {
		return errors.New(fmt.Sprintf("bufLengthErr: %d ,packetBsLne: %d", writeBuf.Len(), packetBytesLength))
	}
	return nil
}
func Test_Connect(t *testing.T) {
	var packetBytesLength int
	bs, err := getPacketBytes(enmu.CONNECT)
	if err != nil {
		t.Fatal(err.Error())
	}
	packetBytesLength = len(bs)
	cp, err := mqtt_packet.ReadOnce(bytes.NewBuffer(bs))
	if err != nil {
		t.Fatal(err.Error())
	}
	p, ok := cp.(*mqtt_packet.ConnectPacket)
	if !ok {
		t.Fatal("is not connect packet")
	}
	if p.ClientIdentifier != "clientid/1" {
		t.Fatal("clientId error: ", p.ClientIdentifier)
	}
	if p.Username != "username/1" {
		t.Fatal("username error: ", p.Username)
	}
	if string(p.Password) != "password" {
		t.Fatal("password error: ", string(p.Password))
	}
	err = writePacket(p, packetBytesLength)
	if err != nil {
		t.Fatal(err.Error())
	}
}
