package packets

//go:generate go run ../../cmd/gen --type=LoginFinished --id=0x02
type LoginFinished struct {
	UUID       UUID
	Name       String
	Properties VarInt
	Strict     Boolean
}
