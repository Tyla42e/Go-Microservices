package models

import (
	"errors"
	"slices"
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

func UpdateComment(c Comment, postId string) error {

	post, ok := posts[postId]

	if !ok {
		return errors.New("Comment not found")
	}

	//comments := post.Comments
	idx := slices.IndexFunc(post.Comments, func(comment Comment) bool { return comment.ID == c.ID })

	if idx == -1 {
		return errors.New("Comment not found")
	}

	post.Comments[idx] = c
	posts[postId] = post
	return nil
}
