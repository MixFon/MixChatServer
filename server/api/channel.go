package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
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

// Обарботка запросов по каналам
func workChannels(w http.ResponseWriter, r *http.Request) {
	fmt.Println("workChannels")
	if r.Method == http.MethodGet {
		getAllChannels(w, r)
	} else if r.Method == http.MethodPost {
		addNewChannel(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

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

	jsChannel, err := json.Marshal(newChannel)
	if err != nil {
		fmt.Println("Error get Channels")
		log.Fatalln("unable marshal to json")
	}
	// Устанавливаем заголовки HTTP для JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем данные в ответ на запрос
	w.Write(jsChannel)
	channels = append(channels, newChannel)
	fmt.Println("newChannel:", newChannel)
	//fmt.Fprintf(w, "Структура Channel успешно сохранена: %+v", newChannel)

}

func StartServer() {
	http.HandleFunc("/channels", workChannels)
	err := http.ListenAndServe(":8080", nil) // устанавливаем порт веб-сервера
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Start server!")
}
