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
	cfg := mysql.Config{
		User:      os.Getenv("DBUSER"),
		Passwd:    os.Getenv("DBPASS"),
		Net:       "tcp",
		Addr:      "127.0.0.1",
		DBName:    "channel",
		ParseTime: true,
	}
	var err error
	// Получаем дескриптор базы данных.
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
	rows, err := dbService.db.Query("SELECT id, name, logoURL, lastMessage, lastActivity FROM Channel")
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса к DB: %w", err)
	}
	defer rows.Close()

	var channels []models.Channel

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
	return channels, nil
}

// Достает канал из базы данных
func getChannelFromRows(rows *sql.Rows) (models.Channel, error) {
	var channel models.Channel

	var lastActivity sql.NullTime
	err := rows.Scan(&channel.Id, &channel.Name, &channel.LogoURL, &channel.LastMessage, &lastActivity)
	if err != nil {
		return models.Channel{}, fmt.Errorf("ошибка чтения структуры channel из DB: %w", err)
	}
	if lastActivity.Valid {
		channel.LastActivity = &lastActivity.Time
	}
	return channel, nil
}

func (dbService *DBService) GetChannel(channelID string) (models.Channel, error) {
	var channel models.Channel
	deleteQuery := "DELETE FROM channel WHERE id = ?"
	rows, err := dbService.db.Query(deleteQuery, channelID)
	defer rows.Close()
	if rows.Next() {
		channel, rows_err := getChannelFromRows(rows)
		if rows_err != nil {
			return channel, fmt.Errorf("ошибка чтения из DB: %w", err)
		}
		return channel, nil
	} else {
		return channel, fmt.Errorf("ошибка итераций по rows: %w", err)
	}
}

// Добавление сообщения по id канала
func (dbService *DBService) AddMessage(channelID string, message models.Message) error {
	insertQuery := "INSERT INTO message (id, text, userID, userName, date, channelID) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := dbService.db.Exec(insertQuery, message.ID, message.Text, message.UserID, message.UserName, message.Date, channelID)
	if err != nil {
		return fmt.Errorf("ошибка добавления сообщения в DB: %w", err)
	}

	updateQuery := "UPDATE channel SET lastMessage = ?, lastActivity = ? WHERE id = ?"
	_, err = dbService.db.Exec(updateQuery, message.Text, message.Date, channelID)
	if err != nil {
		return fmt.Errorf("ошибка обновления даных карала: %w", err)
	}
	return nil
}

// Возвращаем все сообщения по id канала
func (dbService *DBService) GetAllMessages(channelID string) ([]models.Message, error) {
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

	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.ID, &message.Text, &message.UserID, &message.UserName, &message.Date)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}
	// Проверка наличия ошибок после перебора результатов
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("найдена ошибка при переборе результатов: %w", err)
	}
	return messages, nil
}

// Возвращаем все сообщения по id канала
func (dbService *DBService) DelegeChannel(channelID string) error {
	deleteQuery := "DELETE FROM channel WHERE id = ?"
	_, err := dbService.db.Exec(deleteQuery, channelID)
	if err != nil {
		return fmt.Errorf("ошибка удаления карана: %w", err)
	}
	return nil
}
