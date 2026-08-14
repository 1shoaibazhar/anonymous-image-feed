package api

type wsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
