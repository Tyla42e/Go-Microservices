package models

import (
	"errors"
	"slices"
	"time"
)

type Comment struct {
	ID      string    `json:"id"`
	Content string    `binding:"required" json:"content"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}

var comments = make(map[string][]Comment)

func (c Comment) Save(postId string) {

	comments[postId] = append(comments[postId], c)
}

func GetAllComments(postId string) []Comment {

	comments, ok := comments[postId]
	if !ok {
		comments = []Comment{}
	}
	return comments
}

func UpdateComment(c Comment, postId string) error {

	comments, ok := comments[postId]
	if !ok {
		return errors.New("Comment not found")
	}

	idx := slices.IndexFunc(comments, func(comment Comment) bool { return comment.ID == c.ID })

	if idx == -1 {
		return errors.New("Comment not found")
	}

	comments[idx] = c

	return nil
}
