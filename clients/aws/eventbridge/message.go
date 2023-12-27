package eventbridge

// Message presents message event payload
type Message struct {
	// Event Source
	EventSource string
	// Event Detail Type
	EventDetailType string
	// Custom event data
	Detail *DetailEvent
}

type DetailEvent struct {
	EventID string      `json:"eventId"`
	Payload interface{} `json:"payload"`
}
