package main

import (
	"net/http"
	"os"
	"strings"

	"example.com/eventtypes"
	"example.com/moderation/models"
	"example.com/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func main() {

	file, err := os.OpenFile(
		"/var/log/moderation.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	//Sgin.DefaultWriter = io.MultiWriter(file)
	//logger = zerolog.New(file).With().Caller().Timestamp().Logger()
	logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	server := gin.Default()

	server.Use(cors.Default())
	server.POST("/events", handleEvent)
	server.Run(":4003")
}

func handleEvent(context *gin.Context) {
	var event models.CommentEvent

	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
		return
	}
	logger.Info().Msgf("Recieved Event: %+v", event.EventType)
	if event.EventType == eventtypes.CommentCreated.String() {
		if strings.Contains(event.Payload.Content, "orange") {
			event.Payload.Status = "rejected"
		} else {
			event.Payload.Status = "approved"
		}

		event.EventType = eventtypes.CommentModerated.String()
		logger.Info().Msgf("Sending event: %+v\n", event)
		req, err := utils.CreateHTTPRequest("POST", "http://eventbus-srv", "4005", "events", event)

		if err != nil {
			logger.Error().Err(err).Msg("Error Creating Request")
		} else {
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				if res != nil {
					logger.Error().Err(err).Msg(res.Status)
				} else {
					logger.Error().Err(err).Msg("Unable to connect to http://eventbus-srv:4005/events")
				}
			} else {
				logger.Info().Msg(res.Status)
			}
		}
	}
	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}
