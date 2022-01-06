package common

type Request struct {
	ID      string      `json:"id"`
	Args    []string    `json:"args"`
	Data    interface{} `json:"data"`
	Command string      `json:"command"`
}
