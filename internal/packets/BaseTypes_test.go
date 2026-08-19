package packets

import (
	"bytes"
	"testing"
)

func TestVarInt(t *testing.T) {
	testCases := []struct {
		name string
		raw  []byte
		enc  VarInt
	}{
		{"zero", []byte{0x00}, VarInt(0)},
		{"one", []byte{0x01}, VarInt(1)},
		{"max single byte", []byte{0x7f}, VarInt(127)},
		{"min two bytes", []byte{0x80, 0x01}, VarInt(128)},
		{"minus one", []byte{0xff, 0xff, 0xff, 0xff, 0x0f}, VarInt(-1)},
		{"max", []byte{0xff, 0xff, 0xff, 0xff, 0x07}, VarInt(2147483647)},
		{"min", []byte{0x80, 0x80, 0x80, 0x80, 0x08}, VarInt(-2147483648)},
	}

	t.Run("serialise", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got, _ := tc.enc.Serialize()
				if !bytes.Equal(got, tc.raw) {
					t.Errorf("got: %v, want: %v", got, tc.raw)
				}
			})
		}
	})

	t.Run("deserialize", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got, _, _ := DeserializeVarInt(bytes.NewReader(tc.raw))
				if got != tc.enc {
					t.Errorf("got: %v, want: %v", got, tc.enc)
				}
			})
		}
	})

	t.Run("round trip", func(t *testing.T) {
		value := VarInt(1234)
		serialized, _ := value.Serialize()
		deserialized, _, _ := DeserializeVarInt(bytes.NewReader(serialized))
		if deserialized != value {
			t.Errorf("deserialized = %v; want %v", deserialized, value)
		}
	})
}
