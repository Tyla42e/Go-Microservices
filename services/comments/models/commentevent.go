package models

type CommentEvent struct {
	EventType string  `json:"type"`
	PostId    string  `json:"postid"`
	Payload   Comment `json:"data"`
}
