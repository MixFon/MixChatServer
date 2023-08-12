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
	newChannel.LastMessage = nil
	newChannel.LastActivity = nil
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

	err := api.database.DelegeChannel(channelID)
	if err != nil {
		http.Error(w, "ошибка удаления канала", http.StatusBadRequest)
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
		http.Error(w, "Канал не найден", http.StatusBadRequest)
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
	currentTime := time.Now()

	// Форматирование в стандарт ISO 8601
	iso8601Format := "2006-01-02T15:04:05-07:00"
	iso8601Time := currentTime.Format(iso8601Format)
	fmt.Println("ISO 8601 формат:", iso8601Time)

	parsedTime, err := time.Parse(iso8601Format, iso8601Time)
	if err != nil {
		fmt.Println("Ошибка при разборе времени:", err)
		return
	}

	newMessage.Date = parsedTime

	jsMessage, err := json.Marshal(newMessage)
	if err != nil {
		fmt.Println("Err %w", err)
		log.Fatalln("unable marshal to json")
	}

	err = api.database.AddMessage(channelID, newMessage)
	if err != nil {
		http.Error(w, "Ошибка добавления сообщеня в базу данных", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(jsMessage))

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
		fmt.Println("!!!Err %w", err)
		log.Fatalln("unable marshal to json!!!!")
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
