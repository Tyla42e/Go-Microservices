package models

type Event struct {
	Type    string      `json:"type"`
	PostId  string      `json:"postid"`
	Payload interface{} `json:"data"`
}

var events = []Event{}

func (e Event) Save() {
	events = append(events, e)
}

func GetAllEvents() []Event {
	return events
}
