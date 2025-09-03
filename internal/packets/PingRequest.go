package packets

//go:generate go run ../../cmd/gen --type=PingRequest --id=0x01
type PingRequest struct {
	Timestamp Long
}
