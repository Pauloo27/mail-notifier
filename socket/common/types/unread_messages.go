package types

import (
	"time"
)

type CachedUnreadMessages struct {
	Messages []*CachedMailMessage `json:"messages"`
	FechedAt time.Time            `json:"fetched_at"`
}
