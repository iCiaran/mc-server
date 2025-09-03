package packets

//go:generate go run ../../cmd/gen --type=StatusResponse --id=0x00
type StatusResponse struct {
	Response StatusResponseJson `mc:"json"`
}

type StatusResponseVersion struct {
	Name     string `json:"name"`
	Protocol int32  `json:"protocol"`
}

type StatusResponsePlayers struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

type StatusResponseDescription struct {
	Text string `json:"text"`
}

type StatusResponseJson struct {
	Version           StatusResponseVersion     `json:"version"`
	Players           StatusResponsePlayers     `json:"players"`
	Description       StatusResponseDescription `json:"description"`
	EnforceSecureChat bool                      `json:"enforceSecureChat"`
}
