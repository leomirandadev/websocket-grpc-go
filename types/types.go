package types

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
	Channel string
}
