package packets

//go:generate go run ../../cmd/gen --type=LoginFinished --id=0x02
type LoginFinished struct {
	Profile   GameProfile
	SessionId UUID
}
