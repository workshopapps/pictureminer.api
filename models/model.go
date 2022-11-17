package models

import "time"

type UserStruct struct {
	UserName string     `bson:"userName" json:"userName"`
	Email    string     `bson:"email" json:"email"`
	Password string     `bson:"password" json:"password"`
	Request  []Requests `bson:"request" json:"request"`
}
type Requests struct {
	ImageEncode  string    `json:"imageEncode" bson:"imageEncode"`
	DateCreated  time.Time `json:"dateCreated" bson:"dateCreated"`
	RequestsText string    `json:"requestsText" bson:"requestsText"`
}
