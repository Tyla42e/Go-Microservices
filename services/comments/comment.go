package main

import (
	"net/http"
	"os"
	"time"

	"example.com/comments/models"
	"github.com/Tyla42e/Go-Microservices/eventtypes"
	"github.com/Tyla42e/Go-Microservices/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func main() {

	file, err := os.OpenFile(
		"../logs/comment.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	//gin.DefaultWriter = io.MultiWriter(file)
	logger = zerolog.New(file).With().Caller().Timestamp().Logger()

	server := gin.Default()

	server.Use(cors.Default())

	server.GET("/posts/:id/comments", getComments)
	server.POST("/posts/:id/comments", addComment)
	server.POST("/events", handleEvent)
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

	var event models.CommentEvent
	event.EventType = eventtypes.CommentCreated.String()
	event.PostId = postId
	event.Payload = comment
	event.Payload.Status = "pending"

	req, err := utils.CreateHTTPRequest("POST", "http://localhost", "4005", "events", event)

	if err != nil {
		logger.Error().Err(err).Msg("Error Creating Request")
	} else {
		res, err := utils.DispatchRequest(req)
		if err != nil {
			logger.Error().Err(err).Msg(res.Status)
		} else {
			logger.Info().Msg(res.Status)
		}
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Comment Created!", "comment": comment})
}

func handleEvent(context *gin.Context) {
	var event models.CommentEvent

	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
		return
	}

	if event.EventType == eventtypes.CommentModerated.String() {

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
			return
		}

		event.EventType = eventtypes.CommentUpdated.String()

		err = models.UpdateComment(event.Payload, event.PostId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		req, err := utils.CreateHTTPRequest("POST", "http://localhost", "4005", "events", event)

		if err != nil {
			logger.Error().Err(err).Msg("Error Creating Request")
		} else {
			res, err := utils.DispatchRequest(req)
			if err != nil {
				logger.Error().Err(err).Msg(res.Status)
			} else {
				logger.Info().Msg(res.Status)
			}
		}
	}
	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}
