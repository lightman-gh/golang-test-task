package types

import "time"

type Message struct {
	Sender   string    `json:"sender" binding:"required"`
	Receiver string    `json:"receiver" binding:"required"`
	Message  string    `json:"message" binding:"required"`
	Time     time.Time `json:"time"`
}

type Response struct {
	Error string `json:"error"`
	Data  any    `json:"data"`
}
