package test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/qdmc/mqtt_packet"
	"github.com/qdmc/mqtt_packet/enmu"
	"github.com/qdmc/mqtt_packet/packets"
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
func writePacket(cp mqtt_packet.ControlPacketInterface, packetBytesLength int) (*bytes.Buffer, error) {
	writeBuf := bytes.NewBuffer([]byte{})
	var err error
	defer func() {
		if err != nil {
			print("messageType: ", cp.MessageType())
		}

	}()
	length, err := cp.Write(writeBuf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("packetWriteErr: %s", err.Error()))
	}
	if int(length) != packetBytesLength {
		err = errors.New(fmt.Sprintf("packetWriteLenErr: %d ,packetBsLne: %d", length, packetBytesLength))
		return nil, err
	}
	if writeBuf.Len() != packetBytesLength {
		err = errors.New(fmt.Sprintf("bufLengthErr: %d ,packetBsLne: %d", writeBuf.Len(), packetBytesLength))
		return nil, err
	}
	return writeBuf, nil
}

func Test_readPacket(t *testing.T) {
	var packetBytes []byte
	for i, _ := range hexStrMap {
		bs, err := getPacketBytes(i)
		if err != nil {
			t.Fatal(err.Error())
		}
		packetBytes = append(packetBytes, bs...)
		packetBytesLength := len(bs)
		cp, err := mqtt_packet.ReadOnce(bytes.NewBuffer(bs))
		if err != nil {
			t.Fatal(err.Error())
		}
		if cp.MessageType() != i {
			t.Fatal("messageType :", i, "  ", cp.MessageType())
		}
		buf, err := writePacket(cp, packetBytesLength)
		if err != nil {
			t.Fatal(err.Error())
		}
		if string(buf.Bytes()) != string(bs) {
			println("messageType: ", i)
			println("old: ", string(bs), "  hex: ", hex.EncodeToString(bs))
			println("write: ", string(buf.Bytes()), "  hex: ", hex.EncodeToString(buf.Bytes()))
			t.Fatal("packet error")
		}
	}

	packetBytes = append(packetBytes, 0x00)
	println("allLen: ", len(packetBytes))
	list, lastBs, err := mqtt_packet.ReadStream(packetBytes)
	if err != nil {
		t.Fatal("readStreamErr: ", err.Error())
	}
	if lastBs == nil || len(lastBs) != 1 || lastBs[0] != 0x00 {
		println("listLen: ", len(list))
		t.Fatal("lastBytes is error")

	}
	if len(list) != len(hexStrMap) {
		t.Fatal("packet list is error")
	}
}

func Test_PingReq(t *testing.T) {
	var packetBytesLength int
	bs, err := getPacketBytes(enmu.PINGREQ)
	if err != nil {
		t.Fatal(err.Error())
	}
	packetBytesLength = len(bs)
	cp, err := mqtt_packet.ReadOnce(bytes.NewBuffer(bs))
	if err != nil {
		t.Fatal(err.Error())
	}
	p, ok := cp.(*packets.PingReqPacket)
	if !ok {
		t.Fatal("is not pingReq packet")
	}
	buf, err := writePacket(p, packetBytesLength)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(buf.Bytes()) != string(bs) {
		println("messageType: ", enmu.PINGREQ)
		println("old: ", string(bs))
		println("write: ", string(buf.Bytes()))
		t.Fatal("packet error")
	}
}

func Test_UnSubscribe(t *testing.T) {
	var packetBytesLength int
	bs, err := getPacketBytes(enmu.UNSUBSCRIBE)
	if err != nil {
		t.Fatal(err.Error())
	}
	packetBytesLength = len(bs)
	cp, err := mqtt_packet.ReadOnce(bytes.NewBuffer(bs))
	if err != nil {
		t.Fatal(err.Error())
	}
	p, ok := cp.(*mqtt_packet.UnSubscribePacket)
	if !ok {
		t.Fatal("is not unsubscribe packet")
	}
	if p.Topics != nil {
		for i, topic := range p.Topics {
			println("index: ", i, "  topic: ", topic)
		}
	}
	buf, err := writePacket(p, packetBytesLength)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(buf.Bytes()) != string(bs) {
		println("messageType: ", enmu.UNSUBSCRIBE)
		println("old: ", string(bs))
		println("write: ", string(buf.Bytes()))
		t.Fatal("packet error")
	}
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
	buf, err := writePacket(p, packetBytesLength)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(buf.Bytes()) != string(bs) {
		println("messageType: ", enmu.CONNECT)
		println("old: ", string(bs))
		println("write: ", string(buf.Bytes()))
		t.Fatal("packet error")
	}
}

func Test_PubRel(t *testing.T) {
	var packetBytesLength int
	bs, err := getPacketBytes(enmu.PUBREL)
	println("bs :", hex.EncodeToString(bs))
	if err != nil {
		t.Fatal(err.Error())
	}
	packetBytesLength = len(bs)
	cp, err := mqtt_packet.ReadOnce(bytes.NewBuffer(bs))
	if err != nil {
		t.Fatal(err.Error())
	}
	p, ok := cp.(*mqtt_packet.PubRelPacket)
	if !ok {
		t.Fatal("is not PubRelPacket packet")
	}
	buf, err := writePacket(p, packetBytesLength)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(buf.Bytes()) != string(bs) {
		println("messageType: ", enmu.PUBREL)
		println("old: ", string(bs), "  hex: ", hex.EncodeToString(bs))
		println("write: ", string(buf.Bytes()), "  hex: ", hex.EncodeToString(buf.Bytes()))
		t.Fatal("packet error")
	}
}
