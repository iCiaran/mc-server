package packets

//go:generate go run ../../cmd/gen --type=GameProfileProperty

type GameProfileProperty struct {
	Name      String
	Value     String
	Signature Boolean // TODO: change to *String `mc:"prefixed,optional"`
}
