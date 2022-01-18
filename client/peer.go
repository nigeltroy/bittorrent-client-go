package client

type peerConnection struct {
	Id   string `json:"peer id"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type peer struct {
	choked     bool
	connection peerConnection
	interested bool
}
