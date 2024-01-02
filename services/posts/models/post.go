package models

import "time"

type Post struct {
	ID      string    `json:"id"`
	Title   string    `binding:"required" json:"title"`
	Created time.Time `json:"created"`
}

var posts = []Post{}

func (p Post) Save() {
	posts = append(posts, p)
}

func GetAllPosts() []Post {
	return posts
}
