package main

import (
	"net/http"
	"time"

	"example.com/comments/models"
	"example.com/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.Use(cors.Default())

	server.GET("/posts/:id/comments", getComments)
	server.POST("/posts/:id/comments", addComment)
	server.Run(":4001")
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
		return
	}

	idStr, err := utils.GenerateID(4)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
		return
	}

	comment.ID = idStr
	comment.Created = time.Now()

	comment.Save(postId)
	context.JSON(http.StatusCreated, gin.H{"message": "Comment Created!", "comment": comment})
}
