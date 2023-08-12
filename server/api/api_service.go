package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var channels = []models.Channel{}

type DatabaseInterfase interface {
	StartDatabase() error
	AddChannel(channel models.Channel) error
	GetChannels() ([]models.Channel, error)
	AddMessage(channelID string, message models.Message) error
	GetAllMessages(channelID string) ([]models.Message, error)
}

type APIService struct {
	database DatabaseInterfase
}

func (api *APIService) SetDatabase(database DatabaseInterfase) {
	api.database = database
}

// Возвращает список каналов
func (api APIService) GetAllChannels(w http.ResponseWriter, r *http.Request) {
	dbChannels, err := api.database.GetChannels()
	if err != nil {
		fmt.Println("Error get Channels")
		log.Fatalln("unable marshal to json")
	}
	if len(dbChannels) == 0 {
		sendEmptySliceChannels(w, r)
	} else {
		sendSliceChannels(dbChannels, w, r)
	}
}

// Добавление нового канала, присваивает каналу id
func (api APIService) AddNewChannel(w http.ResponseWriter, r *http.Request) {
	// Прочитать тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusInternalServerError)
		return
	}

	var newChannel models.Channel

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
	newChannel.LastActivity = nil
	newChannel.LastMessage = nil
	if err = api.database.AddChannel(newChannel); err != nil {
		http.Error(w, "Ошибка сохранения в базу данных", http.StatusBadRequest)
		return
	}
	sendChannel(newChannel, w, r)
}

// Удаление канала по id
func (api APIService) DeleteChannel(w http.ResponseWriter, r *http.Request) {
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
func (api APIService) GetChannel(w http.ResponseWriter, r *http.Request) {
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
func (api APIService) GetMessagesChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelID := vars["channelID"]

	messages, err := api.database.GetAllMessages(channelID)
	if err != nil {
		log.Fatalln("unable marshal to json")
	}

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
func (api APIService) MessageChannel(w http.ResponseWriter, r *http.Request) {
	// Прочитать тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	channelID := vars["channelID"]

	var newMessage models.Message

	// Декодировать JSON в структуру Channel
	err = json.Unmarshal(body, &newMessage)
	if err != nil {
		http.Error(w, "Ошибка при демаршалинге JSON", http.StatusBadRequest)
		return
	}

	messageID := uuid.New().String()
	newMessage.ID = messageID
	newMessage.Date = time.Now()

	jsMessage, err := json.Marshal(newMessage)
	if err != nil {
		log.Fatalln("unable marshal to json")
	}

	err = api.database.AddMessage(channelID, newMessage)
	if err != nil {
		http.Error(w, "Ошибка добавления сообщеня в базу данных", http.StatusInternalServerError)
		return
	}

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

/// Private

// Кодирование канала и его отправка в формате json
func sendChannel(channel models.Channel, w http.ResponseWriter, r *http.Request) {
	jsChannel, err := json.Marshal(channel)
	if err != nil {
		log.Fatalln("unable marshal to json")
	}
	// Устанавливаем заголовки HTTP для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsChannel)
}

// Отаравна слайса каналов
func sendSliceChannels(channels []models.Channel, w http.ResponseWriter, r *http.Request) {
	jsChannels, err := json.Marshal(channels)
	if err != nil {
		fmt.Println("Error Marshaling Channels")
		log.Fatalln("unable marshal to json")
	}
	// Устанавливаем заголовки HTTP для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsChannels)
}

// Отправка пустого слайса каналов
func sendEmptySliceChannels(w http.ResponseWriter, r *http.Request) {
	emptyArray := []models.Channel{}

	// Преобразуем массив в JSON
	responseJSON, err := json.Marshal(emptyArray)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
