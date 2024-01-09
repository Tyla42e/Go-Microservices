package models

import "time"

type Comment struct {
	ID      string    `json:"id"`
	Content string    `binding:"required" json:"content"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}
