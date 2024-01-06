package main

import (
	"net/http"
	"time"

	"example.com/comments/models"
	"example.com/eventtypes"
	"example.com/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
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

	req, err := utils.CreateHTTPRequest("POST", "http://localhost", "4005", "events", event)

	if err != nil {
		log.Error().Err(err).Msg("Error Creatinmg Request")
	} else {
		res, err := utils.DispatchRequest(req)
		if err != nil {
			log.Error().Err(err).Msg(res.Status)
		} else {
			log.Info().Msg(res.Status)
		}
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Comment Created!", "comment": comment})
}

func handleEvent(context *gin.Context) {
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
	event.EventType = eventtypes.CommentUpdated.String()
	event.PostId = postId
	event.Payload = comment

	req, err := utils.CreateHTTPRequest("POST", "http://localhost", "4005", "events", event)

	if err != nil {
		log.Error().Err(err).Msg("Error Creatinmg Request")
	} else {
		res, err := utils.DispatchRequest(req)
		if err != nil {
			log.Error().Err(err).Msg(res.Status)
		} else {
			log.Info().Msg(res.Status)
		}
	}
	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}
