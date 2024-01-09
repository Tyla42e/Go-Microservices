package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"example.com/eventtypes"
	"example.com/posts/models"
	"example.com/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func main() {

	file, err := os.OpenFile(
		"../services.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	gin.DefaultWriter = io.MultiWriter(file)
	logger = zerolog.New(file).With().Caller().Timestamp().Logger()

	server := gin.Default()

	server.Use(cors.Default())
	server.GET("/posts", getPosts)
	server.POST("/posts", addPost)
	server.POST("/events", handleEvent)
	server.Run(":4000")
}

func getPosts(context *gin.Context) {
	posts := models.GetAllPosts()
	context.JSON(http.StatusOK, posts)
}

func addPost(context *gin.Context) {
	var post models.Post
	err := context.ShouldBindJSON(&post)

	logger.Info().Msgf("Post: %+v", post)
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

	var event models.PostEvent
	event.EventType = eventtypes.PostCreated.String()
	event.Payload = post

	logger.Info().Msgf("Sending event: %+v", event)
	req, err := utils.CreateHTTPRequest("POST", "http://localhost", "4005", "events", event)

	if err != nil {
		logger.Error().Err(err).Msg("Error Creatinmg Request")
	} else {
		res, err := utils.DispatchRequest(req)
		if err != nil {
			logger.Error().Err(err).Msg(res.Status)
		} else {
			logger.Info().Msg(res.Status)
		}
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Post Created!", "post": post})
}

func handleEvent(context *gin.Context) {
	var event models.PostEvent

	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
		return
	}
	logger.Info().Msgf("Recieved Event: %+v", event.EventType)
	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}
