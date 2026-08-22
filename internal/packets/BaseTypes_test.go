package packets

import (
	"testing"
)

func TestVarInt(t *testing.T) {
	testPacket(t, []packetTestCase[VarInt]{
		{"zero", []byte{0x00}, VarInt(0)},
		{"one", []byte{0x01}, VarInt(1)},
		{"max single byte", []byte{0x7f}, VarInt(127)},
		{"min two bytes", []byte{0x80, 0x01}, VarInt(128)},
		{"minus one", []byte{0xff, 0xff, 0xff, 0xff, 0x0f}, VarInt(-1)},
		{"max", []byte{0xff, 0xff, 0xff, 0xff, 0x07}, VarInt(2147483647)},
		{"min", []byte{0x80, 0x80, 0x80, 0x80, 0x08}, VarInt(-2147483648)},
	}, DeserializeVarInt, VarInt(1234))
}

func TestUnsignedShort(t *testing.T) {
	testPacket(t, []packetTestCase[UnsignedShort]{
		{"zero", []byte{0x00, 0x00}, UnsignedShort(0)},
		{"one", []byte{0x00, 0x01}, UnsignedShort(1)},
		{"max", []byte{0xff, 0xff}, UnsignedShort(65535)},
	}, DeserializeUnsignedShort, UnsignedShort(1234))
}

func TestLong(t *testing.T) {
	testPacket(t, []packetTestCase[Long]{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Long(0)},
		{"one", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, Long(1)},
		{"min", []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Long(-9223372036854775808)},
		{"max", []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Long(9223372036854775807)},
	}, DeserializeLong, Long(1234))
}
