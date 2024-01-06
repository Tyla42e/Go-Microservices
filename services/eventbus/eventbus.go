package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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
	//forwardRequest("POST", "http://localhost:4003/events", event)
}

func forwardRequest(method string, url string, data string) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))

	if err != nil {
		log.Error().Err(err).Msg("Error Creatinmg Request")
	} else {
		log.Info().Msg("Requested Created")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg(res.Status)
	} else {
		log.Info().Msg(res.Status)
	}
}
