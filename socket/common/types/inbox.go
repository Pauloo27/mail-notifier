package types

type Inbox struct {
	WebURL   string `json:"web_url"`
	Address  string `json:"address"`
	ID       int    `json:"id"`
	Disabled bool   `json:"disabled"`
}

type Inboxes []*Inbox
