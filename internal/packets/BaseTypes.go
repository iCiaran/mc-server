package packets

import (
	"encoding/binary"
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
	binary.LittleEndian.PutUint16(buffer, uint16(u))
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
