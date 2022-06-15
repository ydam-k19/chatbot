package api

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Route() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/register/admin", RegisterAdmin)
	e.GET("/register/member", RegisterMember)
	e.POST("/meeting/:id", CreateMeeting)
	e.GET("/meeting", GetMeeting)
	e.POST("/meeting/refuse-task", CreateRefuseTask)
	e.PUT("/meeting/refuse-task/:id", UpdateStatusRefuseTask)
	e.POST("/poll", CreatePoll)
	e.PUT("poll/:id", ClosePoll)

	// start server
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

func ConnectDB() {
	// connect to db
	clientOptions := options.Client().ApplyURI("mongodb+srv://ydam:ydam@notes-cluster.iwril.mongodb.net/book?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	DB = client.Database("chatbot")
}

func ConnectBot() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", Bot.Self.UserName)
	Bot.Debug = true
}
