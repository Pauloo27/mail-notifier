package types

import (
	"time"
)

type CachedMailMessage struct {
	ID           string            `json:"id"`
	From         string            `json:"from"`
	Subject      string            `json:"subject"`
	To           []string          `json:"to"`
	Date         time.Time         `json:"date"`
	TextContents map[string][]byte `json:"text_contents"`
	FechedAt     time.Time         `json:"fetched"`
}
