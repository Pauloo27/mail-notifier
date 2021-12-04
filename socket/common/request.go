package common

type Request struct {
	ID      string   `json:"id"`
	Args    []string `json:"args"`
	Command string   `json:"command"`
}
