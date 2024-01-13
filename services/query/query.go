package main

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"

	"example.com/query/models"
	"github.com/Tyla42e/Go-Microservices/eventtypes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var m = map[string]interface{}{
	eventtypes.PostCreated.String():    handlePostCreated,
	eventtypes.CommentCreated.String(): handleCommentCreated,
	eventtypes.CommentUpdated.String(): handleCommentUpdated,
}

var logger zerolog.Logger

func main() {

	file, err := os.OpenFile(
		"../logs/query.log",
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

	server.POST("/events", handleEvent)
	server.GET("/posts", getAllPosts)

	listener, err := net.Listen("tcp", ":4002")
	if err != nil {
		logger.Fatal().Msg("Failed to creat listener")
	}

	go func() {
		err = server.RunListener(listener)
		logger.WithLevel(zerolog.FatalLevel).Err(err).Msg("Failed to start server")
	}()

	processMissedEvents()

	select {}
}

func processMissedEvents() {

	var missedEvents = []models.Event{}

	resp, err := http.Get("http://localhost:4005/events")

	if err != nil {
		logger.Error().Err(err)
		return
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Fatal().Err(err)
			}
			bodyString := string(bodyBytes)

			err = json.Unmarshal([]byte(bodyString), &missedEvents)

			for _, event := range missedEvents {
				logger.Info().Msgf("Processing Event: %+v", event)
				processEvent(event)
			}
		}
	}
}

func handleEvent(context *gin.Context) {
	var event models.Event
	context.ShouldBindJSON(&event)
	logger.Info().Msgf("event: %+v", event)
	logger.Info().Msgf("Type: %s", event.Type)

	processEvent(event)

	context.JSON(http.StatusCreated, gin.H{"message": "OK"})
}

func processEvent(event models.Event) {
	if event.Type != eventtypes.CommentModerated.String() {
		m[event.Type].(func(models.Event))(event)
	}
}

func getAllPosts(context *gin.Context) {
	logger.Info().Msg("getAllPosts")
	data := []*models.Post{}
	posts := models.GetAllPosts()
	for _, post := range posts {
		data = append(data, post)
	}
	context.JSON(http.StatusOK, data)
}

func handlePostCreated(event models.Event) {
	logger.Info().Msg("handlePostCreated")
	logger.Info().Msgf("Event: %+v", event)

	var post models.Post
	jsonData, _ := json.Marshal(event.Payload)
	json.Unmarshal(jsonData, &post)
	post.Comments = []models.Comment{}
	post.Save()
}

func handleCommentCreated(event models.Event) {
	logger.Info().Msg("handleCommentCreated")
	logger.Info().Msgf("Event: %+v", event)

	var comment models.Comment
	jsonData, _ := json.Marshal(event.Payload)
	json.Unmarshal(jsonData, &comment)
	models.AddPost(event.PostId, comment)
}

func handleCommentUpdated(event models.Event) {
	logger.Info().Msg("handleCommentUpdated")
	logger.Info().Msgf("Event: %+v", event)

	var comment models.Comment
	jsonData, _ := json.Marshal(event.Payload)
	json.Unmarshal(jsonData, &comment)

	models.UpdateComment(comment, event.PostId)

}
