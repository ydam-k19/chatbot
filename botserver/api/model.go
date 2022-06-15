package api

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChatID int64              `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Type   string             `json:"type,omitempty" bson:"type,omitempty"`
}

type AdminRegister struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChatID int64              `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Code   string             `json:"code,omitempty" bson:"code,omitempty"`
}

type Task struct {
	MemberID primitive.ObjectID `json:"member_id,omitempty" bson:"member_id,omitempty"`
	Task     string             `json:"task,omitempty" bson:"task,omitempty"`
}
type Venue struct {
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Address   string             `json:"address,omitempty" bson:"address,omitempty"`
	Latitude  float64            `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude float64            `json:"longitude,omitempty" bson:"longitude,omitempty"`
}
type MeetingMessage struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Text      string             `json:"text,omitempty" bson:"text,omitempty"`
	Pin       bool               `json:"pin,omitempty" bson:"pin,omitempty"`
	ShowAdmin bool               `json:"show_admin,omitempty" bson:"show_admin,omitempty"`
	SendTo    []string           `json:"send_to,omitempty" bson:"send_to,omitempty"`
	Venue     *Venue             `json:"venue,omitempty" bson:"venue,omitempty"`
	Tasks     []Task             `json:"tasks,omitempty" bson:"tasks,omitempty"`
	// -------
	AdminID   primitive.ObjectID `json:"admin,omitempty" bson:"admin,omitempty"`
	ChatID    []int64            `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	MessageID []int              `json:"message_id,omitempty" bson:"message_id,omitempty"`
}

type RefuseTask struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChatID    int64              `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	MessageID int                `json:"message_id,omitempty" bson:"message_id,omitempty"`
	MeetingID primitive.ObjectID `json:"meeting_id,omitempty" bson:"meeting_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Reason    string             `json:"reason,omitempty" bson:"reason,omitempty"`
	// -------
	ApproveStatus string `json:"approve_status,omitempty" bson:"approve_status,omitempty"`
}

type Poll struct {
	ID                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Options               []string           `json:"options,omitempty" bson:"options,omitempty"`
	Question              string             `json:"question,omitempty" bson:"question,omitempty"`
	SendTo                string             `json:"send_to,omitempty" bson:"send_to,omitempty"`
	AllowsMultipleAnswers bool               `json:"allows_multiple_answers,omitempty" bson:"allows_multiple_answers,omitempty"`
	// -------
	ChatID    int64 `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	MessageID int   `json:"message_id,omitempty" bson:"message_id,omitempty"`
	// Votes     []int              `json:"votes,omitempty" bson:"votes,omitempty"`
	// Status    string             `json:"status,omitempty" bson:"status,omitempty"`
}
