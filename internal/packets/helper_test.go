package packets

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type packetType interface {
	Serialize() ([]byte, error)
}

type packetTestCase[T any] struct {
	name string
	raw  []byte
	enc  T
}

func testPacket[T packetType](
	t *testing.T,
	cases []packetTestCase[T],
	deserialize func(reader io.Reader) (T, int, error),
	roundTrip T,
) {
	t.Helper()

	t.Run("serialise", func(t *testing.T) {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, _ := tc.enc.Serialize()
				if !bytes.Equal(got, tc.raw) {
					t.Errorf("got: %v, want: %v", got, tc.raw)
				}
			})
		}
	})

	t.Run("deserialize", func(t *testing.T) {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, _, _ := deserialize(bytes.NewReader(tc.raw))
				if !reflect.DeepEqual(got, tc.enc) {
					t.Errorf("got: %v, want: %v", got, tc.enc)
				}
			})
		}
	})

	t.Run("round trip", func(t *testing.T) {
		serialized, _ := roundTrip.Serialize()
		deserialized, _, _ := deserialize(bytes.NewReader(serialized))
		if !reflect.DeepEqual(deserialized, roundTrip) {
			t.Errorf("deserialized = %v; want %v", deserialized, roundTrip)
		}
	})
}
