package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Bot *tgbotapi.BotAPI // bot api
	DB  *mongo.Database  // db
)
