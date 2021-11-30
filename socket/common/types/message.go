package types

import (
	"encoding/json"
	"time"

	"github.com/Pauloo27/mail-notifier/core/provider"
)

type CachedMailMessage struct {
	provider.MailMessage
	FechedAt time.Time
}

func (m *CachedMailMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":      m.GetID(),
		"to":      m.GetTo(),
		"from":    m.GetFrom(),
		"date":    m.GetDate(),
		"subject": m.GetSubject(),
		"fetched": m.FechedAt,
	})
}
