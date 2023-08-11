package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type APIInterface interface {
	GetAllChannels(w http.ResponseWriter, r *http.Request)
	AddNewChannel(w http.ResponseWriter, r *http.Request)
	DeleteChannel(w http.ResponseWriter, r *http.Request)
	GetChannel(w http.ResponseWriter, r *http.Request)
	GetMessagesChannel(w http.ResponseWriter, r *http.Request)
	MessageChannel(w http.ResponseWriter, r *http.Request)
}

func StartServer(api APIInterface) error {
	r := mux.NewRouter()
	r.HandleFunc("/channels", api.GetAllChannels).Methods("Get")
	r.HandleFunc("/channels", api.AddNewChannel).Methods("Post")
	r.HandleFunc("/channels/{channelID}", api.DeleteChannel).Methods("Delete")
	r.HandleFunc("/channels/{channelID}", api.GetChannel).Methods("Get")
	r.HandleFunc("/channels/{channelID}/messages", api.GetMessagesChannel).Methods("Get")
	r.HandleFunc("/channels/{channelID}/messages", api.MessageChannel).Methods("Post")

	http.Handle("/", r)
	fmt.Println("Start server!")
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		return fmt.Errorf("ошибка запуска сервера: %w", err)
	}
	return nil
}
