package models

import (
	"time"
)

type Post struct {
	ID       string    `json:"id"`
	Title    string    `binding:"required" json:"title"`
	Created  time.Time `json:"created"`
	Comments []Comment `json:"comments"`
}

var posts = make(map[string]*Post)

func (p *Post) Save() {
	posts[p.ID] = p
}

func (p Post) UpdateComment(c Comment) {

}

func GetAllPosts() map[string]*Post {
	return posts
}

func AddPost(postId string, c Comment) {
	p := posts[postId]
	p.Comments = append(p.Comments, c)
	posts[postId] = p
}
