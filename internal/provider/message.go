package provider

import "time"

type MailMessage interface {
	GetID() string
	GetSubject() string
	GetTo() []string
	GetFrom() string
	GetDate() time.Time
}
