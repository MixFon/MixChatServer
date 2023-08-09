package main

import (
	"fmt"
	"log"
	"server/api"
	"server/sql_db"
)

func main() {
	fmt.Println("Server start!")
	db := sql_db.DBService{}
	db_err := db.StartDatabase()
	if db_err != nil {
		log.Fatal("Error connect DB: ", db_err)
	}
	api.StartServer(db)
}
