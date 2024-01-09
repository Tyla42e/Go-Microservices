package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	server.POST("/events", handleEvents)
	server.Run(":4005")
}

func handleEvents(context *gin.Context) {

	var result map[string]any
	var data, err = context.GetRawData()
	json.Unmarshal([]byte(data), &result)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
		return
	}

	forwardRequest("POST", "http://localhost:4000/events", string(data[:]))
	forwardRequest("POST", "http://localhost:4001/events", string(data[:]))
	forwardRequest("POST", "http://localhost:4002/events", string(data[:]))
	forwardRequest("POST", "http://localhost:4003/events", string(data[:]))

	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}

func forwardRequest(method string, url string, data string) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))

	logger.Info().Msgf("Sending request to %s", url)
	logger.Info().Msgf("Data:%+v", data)
	if err != nil {
		logger.Error().Err(err).Msg("Error Creatinmg Request")
	} else {
		logger.Info().Msg("Requested Created")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg(res.Status)
	} else {
		logger.Info().Msg(res.Status)
	}
}
