package packets

//go:generate go run ../../cmd/gen --type=GameProfile

type GameProfile struct {
	UUID       UUID
	Username   String
	Properties []GameProfileProperty
}
