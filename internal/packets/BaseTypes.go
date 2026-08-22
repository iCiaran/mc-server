package packets

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
)

type VarInt int32

func DeserializeVarInt(reader io.Reader) (VarInt, int, error) {
	var value int
	position := 0

	for {
		var current byte
		err := binary.Read(reader, binary.BigEndian, &current)
		if err != nil {
			return 0, position + 1, err
		}

		value |= int(current&0x7f) << (position * 7)
		position++

		if position > 5 {
			return 0, position, errors.New("VarInt is too big")
		}

		if (current & 0x80) == 0 {
			break
		}
	}

	return VarInt(value), position, nil
}

func (v VarInt) Serialize() ([]byte, error) {
	buffer := make([]byte, 0)
	for {
		if v & ^0x7f == 0 {
			return append(buffer, byte(v)), nil
		}

		buffer = append(buffer, byte((v&0x7f)|0x80))

		v = VarInt(uint32(v) >> 7)
	}
}

type UnsignedShort uint16

func DeserializeUnsignedShort(reader io.Reader) (UnsignedShort, int, error) {
	data := make([]byte, 2)
	_, err := reader.Read(data)
	if err != nil {
		return 0, 0, err
	}

	return UnsignedShort(binary.BigEndian.Uint16(data)), 0, nil
}

func (u UnsignedShort) Serialize() ([]byte, error) {
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, uint16(u))
	return buffer, nil
}

type Long int64

func DeserializeLong(reader io.Reader) (Long, int, error) {
	data := make([]byte, 8)
	_, err := reader.Read(data)
	if err != nil {
		return 0, 0, err
	}

	return Long(binary.BigEndian.Uint64(data)), 0, nil
}

func (l Long) Serialize() ([]byte, error) {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, uint64(l))
	return buffer, nil
}

type String string

func DeserializeString(reader io.Reader) (String, int, error) {
	stringLength, _, err := DeserializeVarInt(reader)
	if err != nil {
		return "", 0, err
	}

	textBytes := make([]byte, stringLength)
	_, err = reader.Read(textBytes)
	if err != nil {
		return "", 0, err
	}

	return String(textBytes), int(stringLength), nil
}

func (s String) Serialize() ([]byte, error) {
	buffer := make([]byte, 0)

	length, err := VarInt(len(s)).Serialize()
	if err != nil {
		return nil, err
	}

	buffer = append(buffer, length...)
	buffer = append(buffer, []byte(s)...)

	return buffer, nil
}

type Boolean bool

func DeserializeBoolean(reader io.Reader) (Boolean, int, error) {
	buffer := make([]byte, 1)
	_, err := reader.Read(buffer)
	if err != nil {
		return false, 0, err
	}

	return buffer[0] == 1, 0, nil
}

func (b Boolean) Serialize() ([]byte, error) {
	buffer := make([]byte, 1)
	if b {
		buffer[0] = 1
	} else {
		buffer[0] = 0
	}
	return buffer, nil
}

type UUID [16]byte

func DeserializeUUID(reader io.Reader) (UUID, int, error) {
	data := make([]byte, 16)
	_, err := reader.Read(data)
	if err != nil {
		return UUID{}, 0, err
	}
	return UUID(data), 0, nil
}

func (u UUID) Serialize() ([]byte, error) {
	return u[:], nil
}

func (u UUID) String() string {
	var dst [36]byte
	hex.Encode(dst[0:8], u[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], u[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], u[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], u[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], u[10:16])
	return string(dst[:])
}
