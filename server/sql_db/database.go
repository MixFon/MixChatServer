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

func (dbService *DBService) StartDatabase() error {
	// Получаем свойства соединения.
	cfg := mysql.Config{
		User:      os.Getenv("DBUSER"),
		Passwd:    os.Getenv("DBPASS"),
		Net:       "tcp",
		Addr:      "127.0.0.1",
		DBName:    "channel",
		ParseTime: true,
	}
	// Получаем дескриптор базы данных.
	var err error
	dbService.db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("ошибка подключения к mysql: %w", err)
	}

	pingErr := dbService.db.Ping()
	if pingErr != nil {
		return fmt.Errorf("ошибка ping: %w", pingErr)
	}
	fmt.Println("Database connected!")
	return nil
}

// Добавление канала channel
func (dbService *DBService) AddChannel(channel models.Channel) error {
	insertQuery := "INSERT INTO channel (id, name, logoURL) VALUES (?, ?, ?)"
	_, err := dbService.db.Exec(insertQuery, channel.Id, channel.Name, channel.LogoURL)
	if err != nil {
		return fmt.Errorf("ошибка добавления канала: %v", err)
	}
	return nil
}

// Возврат всех каналов
func (dbService DBService) GetAllChannels() ([]models.Channel, error) {
	// Выполнение запроса
	//rows, err := dbService.db.Query("select * from channel")
	rows, err := dbService.db.Query("SELECT id, name, logoURL, lastMessage, lastActivity FROM Channel")
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса к DB: %w", err)
	}
	defer rows.Close()

	// Срез для хранения объектов Channel
	var channels []models.Channel

	// Перебор результатов запроса и создание объектов Channel
	for rows.Next() {
		channel, rows_err := getChannelFromRows(rows)
		if rows_err != nil {
			return nil, fmt.Errorf("ошибка чтения из DB: %w", err)
		}
		channels = append(channels, channel)
	}

	// Проверка наличия ошибок после перебора результатов
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("найдена ошибка при переборе результатов: %w", err)
	}

	// Вывод полученных данных
	// TODO: Убрать
	for _, channel := range channels {
		fmt.Printf("Channel ID: %s\n", channel.Id)
		fmt.Printf("Name: %s\n", channel.Name)
		if channel.LogoURL != nil {
			fmt.Printf("Logo URL: %s\n", *channel.LogoURL)
		}
		/*
			if channel.LastMessage != nil {
				fmt.Printf("Last Message: %s\n", channel.LastMessage)
			}
		*/
		fmt.Printf("Last Activity: %s\n", channel.LastActivity)
		fmt.Println("------------------------")
	}
	return channels, nil
}

// Достает канал из базы данных
func getChannelFromRows(rows *sql.Rows) (models.Channel, error) {
	var channel models.Channel

	var lastActivity sql.NullTime
	err := rows.Scan(&channel.Id, &channel.Name, &channel.LogoURL, &channel.LastMessage, &lastActivity)
	if err != nil {
		fmt.Println("Error!!: %w", err)
		return models.Channel{}, fmt.Errorf("ошибка чтения структуры channel из DB: %w", err)
	}
	if lastActivity.Valid {
		channel.LastActivity = &lastActivity.Time
	}
	return channel, nil
}

func (dbService *DBService) GetChannel(channelID string) (models.Channel, error) {

	var channel models.Channel
	// Подготовка SQL-запроса на удаление
	deleteQuery := "DELETE FROM channel WHERE id = ?"
	rows, err := dbService.db.Query(deleteQuery, channelID)
	defer rows.Close()
	// Перебор результатов запроса и создание объектов Channel
	if rows.Next() {
		channel, rows_err := getChannelFromRows(rows)
		if rows_err != nil {
			return channel, fmt.Errorf("ошибка чтения из DB: %w", err)
		}
		fmt.Printf("Channel with ID %s deleted successfully!\n", channelID)
		return channel, nil
	} else {
		return channel, fmt.Errorf("ошибка итераций по rows: %w", err)
	}
}

// Добавление сообщения по id канала
func (dbService *DBService) AddMessage(channelID string, message models.Message) error {
	// Подготовка SQL-запроса на вставку
	insertQuery := "INSERT INTO message (id, text, userID, userName, date, channelID) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := dbService.db.Exec(insertQuery, message.ID, message.Text, message.UserID, message.UserName, message.Date, channelID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return fmt.Errorf("ошибка добавления сообщения в DB: %w", err)
	}

	fmt.Println("Message added successfully!")
	// Подготовка SQL-запроса на обновление
	updateQuery := "UPDATE channel SET lastMessage = ?, lastActivity = ? WHERE id = ?"
	_, err = dbService.db.Exec(updateQuery, message.Text, message.Date, channelID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return fmt.Errorf("ошибка обновления даных карала: %w", err)
	}

	fmt.Printf("Channel with ID %s updated successfully!\n", channelID)

	return nil
}

// Возвращаем все сообщения по id канала
func (dbService *DBService) GetAllMessages(channelID string) ([]models.Message, error) {
	// SQL-запрос с JOIN для получения сообщений, принадлежащих каналу
	query := `
		SELECT m.id, m.text, m.userID, m.userName, m.date
		FROM message m
		INNER JOIN channel c ON m.channelID = c.id
		WHERE c.id = ?
		ORDER BY m.date ASC
	`
	rows, err := dbService.db.Query(query, channelID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var messages []models.Message

	// Перебор результатов запроса и создание объектов Message
	for rows.Next() {
		//var date sql.NullTime
		var message models.Message
		err := rows.Scan(&message.ID, &message.Text, &message.UserID, &message.UserName, &message.Date)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}

	// Вывод полученных сообщений
	for _, message := range messages {
		fmt.Printf("Message ID: %s\n", message.ID)
		fmt.Printf("Text: %s\n", message.Text)
		fmt.Printf("User ID: %s\n", message.UserID)
		fmt.Printf("User Name: %s\n", message.UserName)
		fmt.Printf("Date: %s\n", message.Date.String())
		fmt.Println("########################")
	}
	return messages, nil
}

// Возвращаем все сообщения по id канала
func (dbService *DBService) DelegeChannel(channelID string) error {
	// Подготовка SQL-запроса на удаление
	deleteQuery := "DELETE FROM channel WHERE id = ?"
	_, err := dbService.db.Exec(deleteQuery, channelID)
	if err != nil {
		return fmt.Errorf("ошибка удаления карана: %w", err)
	}

	fmt.Printf("Channel with ID %s deleted successfully!\n", channelID)
	return nil
}
