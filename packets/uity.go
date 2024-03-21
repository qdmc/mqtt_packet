package packets

import (
	"encoding/binary"
	"errors"
	"io"
)

func readByte(r io.Reader) (byte, error) {
	num := make([]byte, 1)
	_, err := io.ReadFull(r, num)
	//_, err := r.Read(num)
	if err != nil {
		return 0, err
	}
	return num[0], nil
}

func readUin16(r io.Reader) (uint16, error) {
	num := make([]byte, 2)
	_, err := io.ReadFull(r, num)
	//_, err := r.Read(num)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(num), nil
}

func readString(r io.Reader) (int, string, error) {
	var readLen int
	bsLen, err := readUin16(r)
	if err != nil {
		return readLen, "", err
	}
	readLen += 2
	bs, err := readBytes(r, uint(bsLen))
	if err != nil {
		return readLen, "", err
	}
	readLen += int(bsLen)
	return readLen, string(bs), nil
}

func readBytes(r io.Reader, bsLen uint) ([]byte, error) {
	bs := make([]byte, bsLen)
	_, err := io.ReadFull(r, bs)
	//_, err := r.Read(bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// encodeVariableByte 打码 变长字节整数
func encodeVariableByte(length int) ([]byte, error) {
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

// decodeVariableByte 解码码 变长字节整数
func decodeVariableByte(r io.Reader) (int, uint32, error) {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	i := 1
	for ; i < 5; i++ {
		_, err := io.ReadFull(r, b)
		if err != nil {
			return 0, 0, err
		}
		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if digit < 127 {
			break
		}
		multiplier += 7
	}

	return i, rLength, nil
}
