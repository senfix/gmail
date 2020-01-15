package model

import "time"

type Message struct {
	GoogleId string
	Date     time.Time
	From     []Address
	To       []Address
	Subject  string
	Message  string
}

type Address struct {
	Name  string
	Email string
}
