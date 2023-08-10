package sql_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"server/models"
	"time"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type DBService struct {
	db *sql.DB
}

func (dbService *DBService) StartDatabase() error {
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

func (dbService *DBService) AddChannel(channel models.Channel) error {
	insertQuery := "INSERT INTO channel (id, name, logoURL) VALUES (?, ?, ?)"
	_, err := dbService.db.Exec(insertQuery, channel.Id, channel.Name, channel.LogoURL)
	if err != nil {
		return fmt.Errorf("error add channel: %v", err)
	}
	return nil
}

func (dbService DBService) GetChannels() ([]models.Channel, error) {
	// Выполнение запроса
	rows, err := dbService.db.Query("select * from channel")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Срез для хранения объектов Channel
	var channels []models.Channel

	// Перебор результатов запроса и создание объектов Channel
	for rows.Next() {
		var channel models.Channel
		var lastMessage sql.NullString
		var lastActivityStr sql.NullString

		err := rows.Scan(&channel.Id, &channel.Name, &channel.LogoURL, &lastMessage, &lastActivityStr)
		if err != nil {
			log.Fatal(err)
		}

		if lastMessage.Valid {
			channel.LastMessage = &lastMessage.String
		}

		if lastActivityStr.Valid {
			lastActivity, err := time.Parse("2006-01-02 15:04:05", lastActivityStr.String)
			if err != nil {
				log.Fatal(err)
			}
			channel.LastActivity = &lastActivity
		}
		channels = append(channels, channel)
	}

	// Проверка наличия ошибок после перебора результатов
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Вывод полученных данных
	for _, channel := range channels {
		fmt.Printf("Channel ID: %s\n", channel.Id)
		fmt.Printf("Name: %s\n", channel.Name)
		fmt.Printf("Logo URL: %s\n", channel.LogoURL)
		if channel.LastMessage != nil {
			fmt.Printf("Last Message: %s\n", channel.LastMessage)
		}
		if channel.LastActivity != nil {
			fmt.Printf("Last Activity: %s\n", channel.LastActivity.String())
		}
		fmt.Println("------------------------")
	}
	return channels, nil
}
