package models

import "time"

type Comment struct {
	ID      string    `json:"id"`
	Content string    `binding:"required" json:"content"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}

// var comments = make(map[string][]Comment)

// func (c Comment) Save(postId string) {

// 	comments[postId] = append(comments[postId], c)
// }

// func GetAllComments(postId string) []Comment {

// 	comments, ok := comments[postId]
// 	if !ok {
// 		comments = []Comment{}
// 	}
// 	return comments
// }
