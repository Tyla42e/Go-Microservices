package main

import (
	"fmt"
	"net/http"
	"time"

	"example.com/posts/models"
	"example.com/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.Use(cors.Default())
	server.GET("/posts", getPosts)
	server.POST("/posts", addPost)
	server.Run(":4000")
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
		return
	}

	idStr, err := utils.GenerateID(4)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not generate ID"})
		return
	}

	post.ID = idStr
	post.Created = time.Now()

	post.Save()
	context.JSON(http.StatusCreated, gin.H{"message": "Post Created!", "post": post})
}
