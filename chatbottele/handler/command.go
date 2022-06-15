package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetCommands(bot *tgbotapi.BotAPI) {
	setCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "dice",
			Description: "random a number between 1 and 6"},
		tgbotapi.BotCommand{
			Command:     "register",
			Description: "join my team"},
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "join my team"},
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "need help?"},
		tgbotapi.BotCommand{
			Command:     "open_datcom",
			Description: "have a lunch?"},
		tgbotapi.BotCommand{
			Command:     "close_datcom",
			Description: "Have a happy lunch!!"},
		tgbotapi.BotCommand{
			Command:     "vietjet",
			Description: "Dời sharing kkk!!"},
	)

	bot.Request(setCommands)
}

func HandleCommands(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		Register(bot, update)
	case "register":
		Register(bot, update)
	case "vietjet":
		Vietjet(bot, update)
	case "open_datcom":
		OpenDatCom(bot, update)
	case "close_datcom":
		CloseDatCom(bot, update)
	case "dice":
		Dice(bot, update)
	case "help":
		Help(bot, update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func Vietjet(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if update.Message.ReplyToMessage == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This command is only used to reply to meeting messages")
		bot.Send(msg)
	} else {
		if update.Message.Chat.ID == CHANNEL_ID || update.Message.Chat.ID == GROUP_ID {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This command is only used to reply to private messages")
			bot.Send(msg)
			return
		}
		chatID := update.Message.Chat.ID
		messageID := update.Message.ReplyToMessage.MessageID
		resp, err := http.Get(SERVER_URL + "/meeting?chatID=" + strconv.FormatInt(chatID, 10) + "&messageID=" + strconv.Itoa(messageID))
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp.StatusCode != http.StatusNotFound {
				bot.Send(tgbotapi.NewMessage(chatID, "Something went wrong, please try again later"))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "This command is only used to reply to meeting messages"))
			}
			return
		}
		defer resp.Body.Close()
		bodyGet, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
			return
		}

		meetingID := struct {
			ID primitive.ObjectID `json:"_id" bson:"_id"`
		}{}
		if err := json.Unmarshal(bodyGet, &meetingID); err != nil {
			log.Panic(err)
			return
		}
		name := update.Message.From.FirstName + " " + update.Message.From.LastName
		name = strings.Replace(name, " ", "%20", -1)
		link := CLIENT_URL + "/meeting/vietjet?meetingID=" + meetingID.ID.Hex() + "&chatID=" + strconv.FormatInt(chatID, 10) + "&messageID=" + strconv.Itoa(messageID) + "&name=" + name

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please click the link below to refuse the task")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Give the reason", link),
			),
		)
		bot.Send(msg)
	}
}
func CloseDatCom(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "DatCom is now closed. Waiting a few minutes for get your order")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyToMessageID = 0
	bot.Send(msg)
}
func OpenDatCom(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Chọn món đêêêê!!")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("cơm cá hú"),
			tgbotapi.NewKeyboardButton("mì quảng"),
			tgbotapi.NewKeyboardButton("cơm gà"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("cơm rang"),
			tgbotapi.NewKeyboardButton("cơm bò"),
			tgbotapi.NewKeyboardButton("cơm thịt kho trứng"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("bánh canh"),
			tgbotapi.NewKeyboardButton("hủ tiếu mực"),
			tgbotapi.NewKeyboardButton("bánh mì bò kho"),
		),
	)
	bot.Send(msg)
}

func Help(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You need help? Please contact the admin"))
	bot.Send(tgbotapi.NewContact(update.Message.Chat.ID, "0838938301", "Admin"))
}

func Register(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// You are now a member of my team!
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to my team! Please choose your role")

	role := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("admin", "admin"),
			tgbotapi.NewInlineKeyboardButtonData("member", "member"),
		),
	)
	msg.ReplyMarkup = role
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func Dice(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	dice := tgbotapi.NewDice(update.Message.Chat.ID)
	bot.Send(dice)
}
