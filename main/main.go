package main

import (
	"log"
	"server/api"
	"server/handler"
	"server/sql_db"
)

func main() {
	db := sql_db.DBService{}

	// Устанавливаем соединение с DB
	db_err := db.StartDatabase()
	if db_err != nil {
		log.Fatal(db_err)
	}
	api := api.APIService{}
	api.SetDatabase(&db)

	// Запускам сервер
	api_err := handler.StartServer(api)
	if api_err != nil {
		log.Fatal(api_err)
	}
}
