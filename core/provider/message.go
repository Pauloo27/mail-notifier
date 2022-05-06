package provider

import "time"

type MailMessage interface {
	GetID() string
	GetSubject() string
	GetTo() []string
	GetTextContents() map[string][]byte
	GetFrom() string
	GetDate() time.Time
}
