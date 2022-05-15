package storage

type ProviderConfig struct {
	Type     string                 `json:"type"`
	Disabled bool                   `json:"disabled"`
	Info     map[string]interface{} `json:"info"`
}

type Config struct {
	Providers []ProviderConfig `json:"providers"`
}
