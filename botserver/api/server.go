package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func RegisterAdmin(c echo.Context) error {
	var admin AdminRegister

	if err := c.Bind(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if admin.Code != ADMIN_CODE {
		Bot.Send(tgbotapi.NewMessage(admin.ChatID, "Invalid code. Failed to register as admin"))
		return c.JSON(http.StatusBadRequest, "Invalid code")
	}

	insertResult, err := DB.Collection("users").InsertOne(context.TODO(), User{ChatID: admin.ChatID, Name: admin.Name, Type: "admin"})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// send message to user
	msg := tgbotapi.NewMessage(int64(admin.ChatID), "You are now an admin of my team! Please join channel/group below")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("My Team Channel", JOINCHANNELLINK),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("My Team Group", JOINGROUPLINK),
		),
	)
	Bot.Send(msg)

	return c.JSON(http.StatusCreated, insertResult)
}

func RegisterMember(c echo.Context) error {
	chatID, err := strconv.Atoi(c.QueryParam("chatID"))
	name := c.QueryParam("name")
	name = strings.Replace(name, "%20", " ", -1)
	if err != nil {
		return c.JSON(400, "chatID is not a number")
	}
	insertResult, err := DB.Collection("users").InsertOne(context.TODO(), User{ChatID: int64(chatID), Type: "member", Name: name})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// send message to user
	msg := tgbotapi.NewMessage(int64(chatID), "You are now a member of my team! Please join channel/group below")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("My Team Channel", JOINCHANNELLINK),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("My Team Group", JOINGROUPLINK),
		),
	)
	Bot.Send(msg)

	return c.JSON(http.StatusCreated, insertResult)
}

func CreateMeeting(c echo.Context) error {
	var meeting MeetingMessage
	if err := c.Bind(&meeting); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// create meeting
	content := "Meeting: " + "<b>" + meeting.Title + "</b>\n\n" + meeting.Text + "\n\n"
	// tasks
	content += "<b>Tasks</b>\n"

	for i, task := range meeting.Tasks {
		// get member info
		member := User{}
		if err := DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": task.MemberID}).Decode(&member); err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		content += fmt.Sprint(i+1) + ". " + member.Name + " - " + task.Task + "\n"
	}



	// get admin info
	adminID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	meeting.AdminID = adminID
	admin := User{}
	if err := DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": adminID}).Decode(&admin); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	if meeting.ShowAdmin {
		content += "\n\n" + "<b>Admin:</b> " + admin.Name
	}
	
	// forward
	messageID := 0
	var chatID int64
	if contains(meeting.SendTo, CHANNEL) != -1 {
		meeting.ChatID = append(meeting.ChatID, CHANNEL_ID)
		msg := tgbotapi.NewMessage(CHANNEL_ID, content)
		msg.ParseMode = tgbotapi.ModeHTML

		message, _ := Bot.Send(msg)
		// pin
		pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              message.Chat.ID,
			MessageID:           message.MessageID,
			DisableNotification: false,
		}
		_, err := Bot.Request(pinChatMessageConfig)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		// update meeting
		meeting.MessageID = append(meeting.MessageID, message.MessageID)
		messageID = message.MessageID
		chatID = message.Chat.ID
		// venue
		if meeting.Venue != nil {
			venue := tgbotapi.NewVenue(CHANNEL_ID, meeting.Venue.Title, meeting.Venue.Address, meeting.Venue.Latitude, meeting.Venue.Longitude)
			Bot.Send(venue)
		}
	}
	if contains(meeting.SendTo, GROUP) != -1 {
		meeting.ChatID = append(meeting.ChatID, GROUP_ID)
		msg := tgbotapi.NewMessage(GROUP_ID, content)
		msg.ParseMode = tgbotapi.ModeHTML

		message, _ := Bot.Send(msg)
		// pin
		pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              message.Chat.ID,
			MessageID:           message.MessageID,
			DisableNotification: false,
		}
		_, err := Bot.Request(pinChatMessageConfig)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		// update meeting
		meeting.MessageID = append(meeting.MessageID, message.MessageID)
		messageID = message.MessageID
		chatID = message.Chat.ID
		// venue
		if meeting.Venue != nil {
			venue := tgbotapi.NewVenue(GROUP_ID, meeting.Venue.Title, meeting.Venue.Address, meeting.Venue.Latitude, meeting.Venue.Longitude)
			Bot.Send(venue)
		}
	}
	// forward to the member assigned to the task
	for _, task := range meeting.Tasks {
		member := User{}
		if err := DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": task.MemberID}).Decode(&member); err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		msg := tgbotapi.NewForward(member.ChatID, chatID, messageID)
		message, _ := Bot.Send(msg)
		meeting.MessageID = append(meeting.MessageID, message.MessageID)
		meeting.ChatID = append(meeting.ChatID, member.ChatID)
	}

	meeting.Text = content
	// insert meeting to DB
	insertResult, err := DB.Collection("meetings").InsertOne(context.TODO(), meeting)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, insertResult)
}

func GetMeeting(c echo.Context) error {
	chatID, err := strconv.Atoi(c.QueryParam("chatID"))
	if err != nil {
		return c.JSON(400, "chatID is not a number")
	}
	messageID, err := strconv.Atoi(c.QueryParam("messageID"))
	if err != nil {
		return c.JSON(400, "messageID is not a number")
	}
	// if err := DB.Collection("meetings").FindOne(context.TODO(), bson.M{"_id": c.Param("id")}).Decode(&meeting);
	// find in chatID & messageID
	arrChatID := []int64{}
	arrMessageID := []int{}
	arrChatID = append(arrChatID, int64(chatID))
	arrMessageID = append(arrMessageID, messageID)

	meetings := []MeetingMessage{}
	cursor, err := DB.Collection("meetings").Find(context.TODO(), bson.M{"chat_id": bson.M{"$in": arrChatID}, "message_id": bson.M{"$in": arrMessageID}})
	if err != nil {
		return c.JSON(400, err.Error())
	}
	if err := cursor.All(context.TODO(), &meetings); err != nil {
		return c.JSON(400, err.Error())
	}

	for _, meeting := range meetings {
		indexChatID := containsInt64(meeting.ChatID, int64(chatID))
		if indexChatID != -1 && meeting.MessageID[indexChatID] == messageID {
			meetingID := struct {
				ID primitive.ObjectID `json:"_id" bson:"_id"`
			}{
				ID: meeting.ID,
			}
			return c.JSON(200, meetingID)
		}
	}
	return c.JSON(http.StatusNotFound, "Not found")
}

func UpdateStatusRefuseTask(c echo.Context) error {
	refuseTaskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	status := struct {
		Status string `json:"status"`
	}{}
	if err := c.Bind(&status); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// get refuse task
	var refuseTask RefuseTask
	if err := DB.Collection("refuse_tasks").FindOne(context.TODO(), bson.M{"_id": refuseTaskID}).Decode(&refuseTask); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	// update status
	if _, err := DB.Collection("refuse_tasks").UpdateOne(context.TODO(), bson.M{"_id": refuseTaskID}, bson.M{"$set": bson.M{"approve_status": status.Status}}); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if status.Status == "accept" {

		// get meeting
		var meeting MeetingMessage
		if err := DB.Collection("meetings").FindOne(context.TODO(), bson.M{"_id": refuseTask.MeetingID}).Decode(&meeting); err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		// get user
		var user User
		if err := DB.Collection("users").FindOne(context.TODO(), bson.M{"chat_id": refuseTask.ChatID}).Decode(&user); err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		// update meeting
		var newMeeting MeetingMessage
		// task
		index := 0
		removeTask := ""
		for i, task := range meeting.Tasks {
			if task.MemberID != user.ID {
				newMeeting.Tasks = append(newMeeting.Tasks, task)
			} else {
				index = i
				removeTask = task.Task
			}
		}
		// chatID
		for i, chatID := range meeting.ChatID {
			if i != index+len(meeting.SendTo) {
				newMeeting.ChatID = append(newMeeting.ChatID, chatID)
			}
		}
		// messageID
		for i, messageID := range meeting.MessageID {
			if i != index+len(meeting.SendTo) {
				newMeeting.MessageID = append(newMeeting.MessageID, messageID)
			}
		}
		// text
		newTask := fmt.Sprint(index+1) + ". " + user.Name + " - " + removeTask + "\n"
		// split text
		text := strings.Split(meeting.Text, newTask)
		newMeeting.Text = text[0] + "<s>" + newTask + "</s>" + text[1] + "\n<i>(new update)</i>"

		// update meeting
		if _, err := DB.Collection("meetings").UpdateOne(context.TODO(), bson.M{"_id": meeting.ID}, bson.M{"$set": bson.M{"text": newMeeting.Text, "tasks": newMeeting.Tasks, "chat_id": newMeeting.ChatID, "message_id": newMeeting.MessageID}}); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// send update to group/channel
		if contains(meeting.SendTo, CHANNEL) != -1 {
			msg := tgbotapi.NewMessage(CHANNEL_ID, newMeeting.Text)
			msg.ParseMode = tgbotapi.ModeHTML

			message, _ := Bot.Send(msg)
			// pin
			pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
				ChatID:              message.Chat.ID,
				MessageID:           message.MessageID,
				DisableNotification: false,
			}
			_, err := Bot.Request(pinChatMessageConfig)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			// venue
			if meeting.Venue != nil {
				venue := tgbotapi.NewVenue(CHANNEL_ID, meeting.Venue.Title, meeting.Venue.Address, meeting.Venue.Latitude, meeting.Venue.Longitude)
				Bot.Send(venue)
			}
		}
		if contains(meeting.SendTo, GROUP) != -1 {
			msg := tgbotapi.NewMessage(GROUP_ID, newMeeting.Text)
			msg.ParseMode = tgbotapi.ModeHTML

			message, _ := Bot.Send(msg)
			// pin
			pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
				ChatID:              message.Chat.ID,
				MessageID:           message.MessageID,
				DisableNotification: false,
			}
			_, err := Bot.Request(pinChatMessageConfig)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			// venue
			if meeting.Venue != nil {
				venue := tgbotapi.NewVenue(GROUP_ID, meeting.Venue.Title, meeting.Venue.Address, meeting.Venue.Latitude, meeting.Venue.Longitude)
				Bot.Send(venue)
			}
		}
		// send message to user
		msg := tgbotapi.NewMessage(refuseTask.ChatID, "Admin accept your request")
		msg.ReplyToMessageID = refuseTask.MessageID
		Bot.Send(msg)

	} else if status.Status == "reject" {
		// send message to user
		msg := tgbotapi.NewMessage(refuseTask.ChatID, "Your request has been rejected")
		msg.ReplyToMessageID = refuseTask.MessageID
		Bot.Send(msg)
	}

	return c.JSON(http.StatusOK, "ok")
}

func CreateRefuseTask(c echo.Context) error {
	var refuseTask RefuseTask
	if err := c.Bind(&refuseTask); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	refuseTask.Name = strings.Replace(refuseTask.Name, "%20", " ", -1)
	// get Meeting
	var meeting MeetingMessage
	if err := DB.Collection("meetings").FindOne(context.TODO(), bson.M{"_id": refuseTask.MeetingID}).Decode(&meeting); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	// get Admin
	var admin User
	if err := DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": meeting.AdminID}).Decode(&admin); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	refuseTask.ApproveStatus = APPROVE_STATUS_PENDING
	// create refuseTask DB
	insertResult, err := DB.Collection("refuse_tasks").InsertOne(context.TODO(), refuseTask)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// forward to admin
	content := refuseTask.Name + " want to refuse task with the reason: " + refuseTask.Reason
	msg, err := Bot.Send(tgbotapi.NewForward(admin.ChatID, refuseTask.ChatID, refuseTask.MessageID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	replyMsg := tgbotapi.NewMessage(admin.ChatID, content)
	replyMsg.ReplyToMessageID = msg.MessageID

	insertID, _ := insertResult.InsertedID.(primitive.ObjectID)

	options := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("accept", "accept RefuseTaskID:"+insertID.Hex()),
			tgbotapi.NewInlineKeyboardButtonData("reject", "reject RefuseTaskID:"+insertID.Hex()),
		),
	)
	replyMsg.ReplyMarkup = &options
	Bot.Send(replyMsg)

	return c.JSON(http.StatusOK, insertResult)
}

func CreatePoll(c echo.Context) error {
	var poll Poll
	if err := c.Bind(&poll); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if poll.SendTo == CHANNEL {
		poll.ChatID = CHANNEL_ID
	} else if poll.SendTo == GROUP {
		poll.ChatID = GROUP_ID
	}

	newPoll := tgbotapi.NewPoll(poll.ChatID, poll.Question, poll.Options...)
	newPoll.AllowsMultipleAnswers = poll.AllowsMultipleAnswers
	msg, _ := Bot.Send(newPoll)

	poll.MessageID = msg.MessageID

	insertResult, err := DB.Collection("polls").InsertOne(context.TODO(), poll)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, insertResult)
}

func ClosePoll(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id")) // click close tren giao dien

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// get poll
	var poll Poll
	if err := DB.Collection("polls").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&poll); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	Bot.StopPoll(tgbotapi.NewStopPoll(poll.ChatID, poll.MessageID))

	return c.JSON(http.StatusOK, "ok")
}
