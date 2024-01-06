package models

type PostEvent struct {
	EventType string `json:"type"`
	Payload   Post   `json:"data"`
}
