package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/query/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var m = map[string]interface{}{
	"PostCreated":    handlePostCreated,
	"CommentCreated": handleCommentCreated,
}

func main() {

	server := gin.Default()

	server.Use(cors.Default())

	server.POST("/events", handleEvent)
	server.GET("/posts", getAllPosts)
	server.Run(":4002")
}

func handleEvent(context *gin.Context) {
	var event models.Event
	context.ShouldBindJSON(&event)
	m[event.Type].(func(models.Event))(event)

	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}

func getAllPosts(context *gin.Context) {

	data := []*models.Post{}
	posts := models.GetAllPosts()
	fmt.Printf("posts: %+v\n", posts)
	for _, v := range posts {
		var post *models.Post
		post = v
		data = append(data, post)
	}
	context.JSON(http.StatusOK, data)
}

func handlePostCreated(event models.Event) {
	var post models.Post
	jsonData, _ := json.Marshal(event.Payload)
	json.Unmarshal(jsonData, &post)
	post.Comments = []models.Comment{}
	post.Save()
}

func handleCommentCreated(event models.Event) {
	var comment models.Comment
	jsonData, _ := json.Marshal(event.Payload)
	json.Unmarshal(jsonData, &comment)
	models.AddPost(event.PostId, comment)
}
