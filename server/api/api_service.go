package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DatabaseInterfase interface {
	StartDatabase() error
	AddChannel(channel models.Channel) error
	GetAllChannels() ([]models.Channel, error)
	AddMessage(channelID string, message models.Message) error
	GetAllMessages(channelID string) ([]models.Message, error)
	DelegeChannel(channelID string) error
	GetChannel(channelID string) (models.Channel, error)
}

type APIService struct {
	database DatabaseInterfase
}

func (api *APIService) SetDatabase(database DatabaseInterfase) {
	api.database = database
}

// Возвращает список каналов
func (api APIService) GetAllChannels(w http.ResponseWriter, r *http.Request) {
	dbChannels, err := api.database.GetAllChannels()
	if err != nil {
		sendError(err, w, r)
		return
	}
	if len(dbChannels) == 0 {
		sendEmptySliceChannels(w, r)
	} else {
		sendSliceChannels(dbChannels, w, r)
	}
}

// Добавление нового канала, присваивает каналу id
func (api APIService) AddNewChannel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(err, w, r)
		return
	}

	var newChannel models.Channel

	err = json.Unmarshal(body, &newChannel)
	if err != nil {
		sendError(err, w, r)
		return
	}

	id := uuid.New()
	newChannel.Id = id.String()
	newChannel.LastMessage = nil
	newChannel.LastActivity = nil
	if err = api.database.AddChannel(newChannel); err != nil {
		sendError(err, w, r)
		return
	}
	sendChannel(newChannel, w, r)
}

// Удаление канала по id
func (api APIService) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelID := vars["channelID"]

	err := api.database.DelegeChannel(channelID)
	if err != nil {
		sendError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

// Возвращает канал по id канала
func (api APIService) GetChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelID := vars["channelID"]

	channel, err := api.database.GetChannel(channelID)
	if err != nil {
		sendError(err, w, r)
		return
	}
	sendChannel(channel, w, r)
}

// Отправляем сообщения канала
func (api APIService) GetMessagesChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelID := vars["channelID"]

	messages, err := api.database.GetAllMessages(channelID)
	if err != nil {
		sendError(err, w, r)
		return
	}

	jsMessages, err := json.Marshal(messages)
	if err != nil {
		sendError(err, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsMessages)
}

// Получение нового сообщения в канал
func (api APIService) MessageChannel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(err, w, r)
		return
	}

	vars := mux.Vars(r)
	channelID := vars["channelID"]

	var newMessage models.Message
	err = json.Unmarshal(body, &newMessage)
	if err != nil {
		sendError(err, w, r)
		return
	}

	messageID := uuid.New().String()
	newMessage.ID = messageID
	currentTime := time.Now()

	// Форматирование в стандарт ISO 8601
	iso8601Format := "2006-01-02T15:04:05-07:00"
	iso8601Time := currentTime.Format(iso8601Format)
	parsedTime, err := time.Parse(iso8601Format, iso8601Time)
	if err != nil {
		sendError(err, w, r)
		return
	}

	newMessage.Date = parsedTime
	jsMessage, err := json.Marshal(newMessage)
	if err != nil {
		sendError(err, w, r)
		return
	}

	err = api.database.AddMessage(channelID, newMessage)
	if err != nil {
		sendError(err, w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsMessage)
}

/// Private

// Кодирование канала и его отправка в формате json
func sendChannel(channel models.Channel, w http.ResponseWriter, r *http.Request) {
	jsChannel, err := json.Marshal(channel)
	if err != nil {
		sendError(err, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsChannel)
}

// Отаравна слайса каналов
func sendSliceChannels(channels []models.Channel, w http.ResponseWriter, r *http.Request) {
	jsChannels, err := json.Marshal(channels)
	if err != nil {
		sendError(err, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsChannels)
}

// Отправка пустого слайса каналов
func sendEmptySliceChannels(w http.ResponseWriter, r *http.Request) {
	emptyArray := []models.Channel{}
	responseJSON, err := json.Marshal(emptyArray)
	if err != nil {
		sendError(err, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// Отправка ошибки
func sendError(err error, w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprint("%w\n", err), http.StatusInternalServerError)
}
