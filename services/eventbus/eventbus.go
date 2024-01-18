package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"example.com/eventbus/models"
	"github.com/rs/zerolog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var logger zerolog.Logger

func main() {

	file, err := os.OpenFile(
		"/var/log/eventbus.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	//gin.DefaultWriter = io.MultiWriter(file)
	//logger = zerolog.New(file).With().Caller().Timestamp().Logger()

	logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	server := gin.Default()

	server.Use(cors.Default())
	server.GET("/", getAllEvents)
	server.POST("/events", handleEvents)
	server.GET("/events", getAllEvents)
	server.Run(":4005")
}

func getAllEvents(context *gin.Context) {
	events := models.GetAllEvents()
	context.JSON(http.StatusOK, events)
}

func handleEvents(context *gin.Context) {

	logger.Info().Msgf("Headers: %+v", context.Request.Header)
	var event models.Event

	var err = context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
		return
	}

	logger.Info().Msgf("Recieved Event: %+v", event)

	event.Save()

	logger.Info().Msg("Forwarding Requests")
	forwardRequest("POST", "http://post-clusterip-srv:4000/events", event)
	forwardRequest("POST", "http://comments-clusterip-srv:4001/events", event)
	forwardRequest("POST", "http://query-clusterip-srv:4002/events", event)
	forwardRequest("POST", "http://moderation-clusterip-srv:4003/events", event)

	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}

func forwardRequest(method string, url string, event models.Event) {

	reqBodyBytes, err := json.MarshalIndent(event, "", "\t")
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBodyBytes))

	logger.Info().Msgf("Sending request to %s", url)
	logger.Info().Msgf("Request Data: %+v", string(reqBodyBytes[:]))

	if err != nil {
		logger.Error().Err(err).Msg("Error Creating Request")
	} else {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			if res != nil {
				logger.Error().Err(err).Msg(res.Status)
			} else {
				logger.Error().Err(err).Msgf("Unable to connect to %s", url)
			}
		} else {
			logger.Info().Msg(res.Status)
		}
	}
}
