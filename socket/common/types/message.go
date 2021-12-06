package types

import (
	"encoding/json"
	"time"
)

type CachedMailMessage struct {
	ID, From, Subject string
	To                []string
	Date              time.Time
	FechedAt          time.Time
}

func (m *CachedMailMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":      m.ID,
		"to":      m.To,
		"from":    m.From,
		"date":    m.Date,
		"subject": m.Subject,
		"fetched": m.FechedAt,
	})
}
