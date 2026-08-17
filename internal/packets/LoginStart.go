package packets

//go:generate go run ../../cmd/gen --type=LoginStart --id=0x00
type LoginStart struct {
	Name String
	UUID UUID
}
