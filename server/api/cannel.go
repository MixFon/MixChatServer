package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

/*
public struct Channel: Codable {
    public let id: String
    public let name: String
    public let logoURL: String?
    public let lastMessage: String?
    public let lastActivity: Date?
}
*/

type Channel struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LogoURL      string `json:"logoURL"`
	LastMessage  string `json:"lastMessage,omitempty"`
	LastActivity string `json:"lastActivity,omitempty"`
}

func getChannels(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	strId := id.String()
	channel := Channel{
		Id:      strId,
		Name:    "Mix",
		LogoURL: "Logo",
	}
	jsChannel, err := json.Marshal(channel)
	if err != nil {
		log.Fatalln("unable marshal to json")
	}
	fmt.Fprintf(w, "Channel %v", string(jsChannel))
}

func StartServer() {
	http.HandleFunc("/channels", getChannels)
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Start server!")
}
