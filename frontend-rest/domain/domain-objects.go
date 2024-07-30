package domain

type Event struct { // In-Clickhouse event
	// Event IDentifier
	EID        string `json:"eid" ch:"eid"`
	ProviderID string `json:"provider_id" ch:"provider_id"`
	SchemaID   string `json:"schema_id" ch:"schema_id"`
	Type       string `json:"type" ch:"type"`
	Data       string `json:"data" ch:"data"`
}
