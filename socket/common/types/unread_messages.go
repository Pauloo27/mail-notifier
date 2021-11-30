package types

import (
	"encoding/json"
	"time"
)

type CachedUnreadMessages struct {
	Messages []*CachedMailMessage
	FechedAt time.Time
}

func (m *CachedUnreadMessages) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"messages": m.Messages,
		"fetched":  m.FechedAt,
	})
}
