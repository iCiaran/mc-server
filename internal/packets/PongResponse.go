package packets

//go:generate go run ../../cmd/gen --type=PongResponse --id=0x01
type PongResponse struct {
	Timestamp Long
}
