package main

import (
	"fmt"
	"net/http"
	"time"

	"example.com/models"
	"example.com/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.Use(cors.Default())

	server.GET("/posts", getPosts)
	server.POST("/posts", addPost)
	server.GET("/posts/:id/comments", getComments)
	server.POST("/posts/:id/comments", addComment)

	server.Run(":8080")
}

func getPosts(context *gin.Context) {
	posts := models.GetAllPosts()
	context.JSON(http.StatusOK, posts)
}

func addPost(context *gin.Context) {
	var post models.Post
	fmt.Println(context.Params)
	err := context.ShouldBindJSON(&post)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
	}

	idStr, err := utils.GenerateID(4)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
	}

	post.ID = idStr
	post.Created = time.Now()

	post.Save()
	context.JSON(http.StatusCreated, gin.H{"message": "Post Created!", "post": post})
}

func getComments(context *gin.Context) {
	postId := context.Param("id")
	posts := models.GetAllComments(postId)
	context.JSON(http.StatusOK, posts)
}

func addComment(context *gin.Context) {
	postId := context.Param("id")
	var comment models.Comment

	err := context.ShouldBindJSON(&comment)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
	}

	idStr, err := utils.GenerateID(4)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
	}

	comment.ID = idStr
	comment.Created = time.Now()

	comment.Save(postId)
	context.JSON(http.StatusCreated, gin.H{"message": "Comment Created!", "comment": comment})
}