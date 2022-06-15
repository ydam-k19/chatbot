package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCallbackQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if strings.Contains(update.CallbackQuery.Data, " RefuseTaskID:") {
		refuseTaskID := strings.Split(update.CallbackQuery.Data, " RefuseTaskID:")[1]
		update.CallbackQuery.Data = strings.Split(update.CallbackQuery.Data, " RefuseTaskID:")[0]
		// send request to update the status of the task
		client := &http.Client{}
		status := struct {
			Status string `json:"status"`
		}{
			Status: update.CallbackQuery.Data,
		}
		json_data, _ := json.Marshal(status)
		req, err := http.NewRequest(http.MethodPut, SERVER_URL+"/meeting/refuse-task/"+refuseTaskID, bytes.NewBuffer(json_data))
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Something went wrong, please try again later"))
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Something went wrong, please try again later"))
		}
	}
	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("\nYou chose option: "+update.CallbackQuery.Data, update.CallbackQuery.Data),
			},
		},
	}

	bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, markup))

	if update.CallbackQuery.Data == "admin" {
		chatID := strconv.FormatInt(update.CallbackQuery.From.ID, 10)
		name := update.CallbackQuery.From.FirstName + " " + update.CallbackQuery.From.LastName
		msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Please verify your information by entering the code in the form below")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("become an admin", CLIENT_URL+"/register/admin?chatID="+chatID+"&name="+name),
			),
		)
		bot.Send(msg)
	} else if update.CallbackQuery.Data == "member" {
		chatID := strconv.FormatInt(update.CallbackQuery.From.ID, 10)
		name := update.CallbackQuery.From.FirstName + " " + update.CallbackQuery.From.LastName
		name = strings.Replace(name, " ", "%20", -1)
		resp, err := http.Get(SERVER_URL + "/register/member?chatID=" + chatID + "&name=" + name)

		if err != nil || resp.StatusCode != http.StatusCreated {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Something went wrong, please try again later"))
		}
	}

}
