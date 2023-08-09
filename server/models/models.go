package models

import "time"

type Channel struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	LogoURL      string     `json:"logoURL"`
	LastMessage  *string    `json:"lastMessage,omitempty"`
	LastActivity *time.Time `json:"lastActivity,omitempty"`
}

type Message struct {
	ID       string    `json:"id"`
	Text     string    `json:"text"`
	UserID   string    `json:"userID"`
	UserName string    `json:"userName"`
	Date     time.Time `json:"date"`
}
