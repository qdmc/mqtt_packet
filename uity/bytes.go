package uity

import (
	"encoding/binary"
	"errors"
	"io"
)

func ReadByte(r io.Reader) (byte, error) {
	num := make([]byte, 1)
	_, err := io.ReadFull(r, num)
	//_, err := r.Read(num)
	if err != nil {
		return 0, err
	}
	return num[0], nil
}

func ReadUin16(r io.Reader) (uint16, error) {
	num := make([]byte, 2)
	_, err := io.ReadFull(r, num)
	//_, err := r.Read(num)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(num), nil
}

func ReadString(r io.Reader) (string, error) {
	bsLen, err := ReadUin16(r)
	if err != nil {
		return "", err
	}
	bs, err := ReadBytes(r, uint(bsLen))
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func ReadBytes(r io.Reader, bsLen uint) ([]byte, error) {
	bs := make([]byte, bsLen)
	_, err := io.ReadFull(r, bs)
	//_, err := r.Read(bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// EncodeVariableByte 打码 变长字节整数
func EncodeVariableByte(length int) ([]byte, error) {
	maxLen := 268435455
	if length > maxLen {
		return nil, errors.New("length is to big")
	}
	if length < 0 {
		return nil, errors.New("length is must >= 0")
	}
	if length == 0 {
		return []byte{}, nil
	}
	var encLength []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encLength = append(encLength, digit)
		if length == 0 {
			break
		}
	}
	return encLength, nil
}

// DecodeVariableByte 解码码 变长字节整数
func DecodeVariableByte(r io.Reader) (int, error) {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	for i := 0; i < 4; i++ {
		_, err := io.ReadFull(r, b)
		if err != nil {
			return 0, err
		}
		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if digit < 127 {
			break
		}
		multiplier += 7
	}
	return int(rLength), nil
}
