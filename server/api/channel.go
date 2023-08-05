package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Channel struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LogoURL      string `json:"logoURL"`
	LastMessage  string `json:"lastMessage,omitempty"`
	LastActivity string `json:"lastActivity,omitempty"`
}

var channels = []Channel{
	{
		Id:      "channel1",
		Name:    "Channel One",
		LogoURL: "https://images.dog.ceo/breeds/entlebucher/n02108000_2212.jpg",
	},
	{
		Id:      "channel2",
		Name:    "Channel Two",
		LogoURL: "https://images.dog.ceo/breeds/pyrenees/n02111500_124.jpg",
	},
}

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

	// Отправляем данные в ответ на запрос
	w.Write(jsChannels)
	fmt.Println("Success get Channels!")
}

// Добавление нового канала, присваивает каналу id
func addNewChannel(w http.ResponseWriter, r *http.Request) {
	// Прочитать тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusInternalServerError)
		return
	}
	// Создать переменную для декодирования JSON
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

	// Отправляем данные в ответ на запрос
	w.Write(jsChannel)
}

func deleteChannel(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение переменной из URL-пути по ключу "id".
	vars := mux.Vars(r)
	userID := vars["id"]

	for i, channel := range channels {
		if channel.Id == userID {
			channels = append(channels[:i], channels[i+1:]...)
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
			return
		}
	}
	http.Error(w, "Канал не найден", http.StatusBadRequest)
}

// Возвращает канад по id канала
func getChannel(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение переменной из URL-пути по ключу "id".
	vars := mux.Vars(r)
	userID := vars["id"]

	for _, channel := range channels {
		if channel.Id == userID {
			sendChannel(channel, w, r)
			return
		}
	}
	http.Error(w, "Канал не найден", http.StatusBadRequest)
}

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/channels", getAllChannels).Methods("Get")
	r.HandleFunc("/channels", addNewChannel).Methods("Post")
	r.HandleFunc("/channels/{id}", deleteChannel).Methods("Delete")
	r.HandleFunc("/channels/{id}", deleteChannel).Methods("Get")

	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Start server!")
}
