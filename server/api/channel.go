package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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

var messageDict = make(map[string][]Message)

var channels = []Channel{}

// Возвращает список каналов
func getAllChannels(w http.ResponseWriter, r *http.Request) {
	jsChannels, err := json.Marshal(channels)
	if err != nil {
		fmt.Println("Error get Channels")
		log.Fatalln("unable marshal to json")
	}
	// Устанавливаем заголовки HTTP для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsChannels)
}

// Добавление нового канала, присваивает каналу id
func addNewChannel(w http.ResponseWriter, r *http.Request) {
	// Прочитать тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusInternalServerError)
		return
	}

	var newChannel Channel

	// Декодировать JSON в структуру Channel
	err = json.Unmarshal(body, &newChannel)
	if err != nil {
		http.Error(w, "Ошибка при демаршалинге JSON", http.StatusBadRequest)
		return
	}

	// В этом моменте у вас есть новая структура Channel в переменной newChannel.
	// Здесь вы можете произвести сохранение новой структуры в базу данных или в другое место по вашему выбору.
	// Например, можно добавить ее в существующий слайс или сохранить в базу данных.
	id := uuid.New()
	newChannel.Id = id.String()
	channels = append(channels, newChannel)
	sendChannel(newChannel, w, r)
}

// Кодирование канала и его отправка в формате json
func sendChannel(channel Channel, w http.ResponseWriter, r *http.Request) {
	jsChannel, err := json.Marshal(channel)
	if err != nil {
		log.Fatalln("unable marshal to json")
	}
	// Устанавливаем заголовки HTTP для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsChannel)
	fmt.Println(string(jsChannel))
}

// Удаление канала по id
func deleteChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelID := vars["channelID"]

	for i, channel := range channels {
		if channel.Id == channelID {
			channels = append(channels[:i], channels[i+1:]...)
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
			return
		}
	}
	http.Error(w, "Канал не найден", http.StatusBadRequest)
}

// Возвращает канал по id канала
func getChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelID := vars["channelID"]

	for _, channel := range channels {
		if channel.Id == channelID {
			sendChannel(channel, w, r)
			return
		}
	}
	http.Error(w, "Канал не найден", http.StatusBadRequest)
}

// Отправляем сообщения канала
func getMessagesChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelID := vars["channelID"]

	messages := messageDict[channelID]

	jsMessages, err := json.Marshal(messages)
	if err != nil {
		log.Fatalln("unable marshal to json")
	}
	// Устанавливаем заголовки HTTP для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsMessages)
}

// Получение нового сообщения в канал
func messageChannel(w http.ResponseWriter, r *http.Request) {
	// Прочитать тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	channelID := vars["channelID"]

	var newMessage Message

	// Декодировать JSON в структуру Channel
	err = json.Unmarshal(body, &newMessage)
	if err != nil {
		http.Error(w, "Ошибка при демаршалинге JSON", http.StatusBadRequest)
		return
	}

	newMessage.ID = channelID
	newMessage.Date = time.Now()

	jsMessage, err := json.Marshal(newMessage)
	if err != nil {
		log.Fatalln("unable marshal to json")
	}

	messageDict[channelID] = append(messageDict[channelID], newMessage)

	for i, channel := range channels {
		if channel.Id == channelID {
			channels[i].LastActivity = &newMessage.Date
			channels[i].LastMessage = &newMessage.Text
			break
		}
	}

	// Устанавливаем заголовки HTTP для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsMessage)
}

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/channels", getAllChannels).Methods("Get")
	r.HandleFunc("/channels", addNewChannel).Methods("Post")
	r.HandleFunc("/channels/{channelID}", deleteChannel).Methods("Delete")
	r.HandleFunc("/channels/{channelID}", getChannel).Methods("Get")
	r.HandleFunc("/channels/{channelID}/messages", getMessagesChannel).Methods("Get")
	r.HandleFunc("/channels/{channelID}/messages", messageChannel).Methods("Post")

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Start server!")
}
