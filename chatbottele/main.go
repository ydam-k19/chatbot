package main

import (
	"log"

	"chatbottele/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(handler.BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}
	handler.SetCommands(bot)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handler.HandleCommands(bot, &update)
			}
		} else if update.CallbackQuery != nil {
			handler.HandleCallbackQuery(bot, &update)
		}
	}
}
