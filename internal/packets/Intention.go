package packets

//go:generate go run ../../cmd/gen --type=Intention --id=0x01
type Intention struct {
	ProtocolVersion VarInt
	ServerAddress   String
	ServerPort      UnsignedShort
	Intent          VarInt
}
