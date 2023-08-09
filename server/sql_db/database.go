package sql_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"server/models"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type DBService struct {
	db *sql.DB
}

func (dbService DBService) StartDatabase() error {
	// Получаем свойства соединения.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1",
		DBName: "channel",
	}
	// Получаем дескриптор базы данных.
	var err error
	dbService.db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return err
	}

	pingErr := dbService.db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return err
	}
	fmt.Println("Database connected!")
	return nil
}

func (dbService DBService) AddChannel(channel models.Channel) error {
	_, err := dbService.db.Exec("INSERT INTO channel (id, name, logoURL, lastMessage, lastActivity) VALUES (?, ?, ?, ?, ?)", channel.Id, channel.Name, channel.LastMessage, channel.LastActivity)
	if err != nil {
		return fmt.Errorf("error add channel: %v", err)
	}
	return nil
}

func (dbService DBService) getChannels() ([]models.Channel, error) {

	return []models.Channel{}, nil
}
