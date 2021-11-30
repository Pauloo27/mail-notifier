package types

import (
	"encoding/json"

	"github.com/Pauloo27/mail-notifier/core/provider"
)

type Inbox struct {
	provider.MailBox
}

type Inboxes []*Inbox

func (i *Inbox) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"address": i.GetAddress(),
	})
}
