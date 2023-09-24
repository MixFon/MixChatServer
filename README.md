# MixChatServer
Данная программа обеспечивает работу мобильного приложения [MixChat](https://github.com/MixFon/MixChat). 

# Документация API

## Методы

| Метод    | Путь                              | Функция-обработчик    | Описание                             |
|-----------|-----------------------------------|-----------------------|--------------------------------------|
| GET       | `/channels`                       | `api.GetAllChannels`  | Получить список всех каналов.       |
| POST      | `/channels`                       | `api.AddNewChannel`   | Добавить новый канал.               |
| DELETE    | `/channels/{channelID}`           | `api.DeleteChannel`   | Удалить канал по его `channelID`.   |
| GET       | `/channels/{channelID}`           | `api.GetChannel`      | Получить информацию о канале.       |
| GET       | `/channels/{channelID}/messages`  | `api.GetMessagesChannel` | Получить сообщения в канале.    |
| POST      | `/channels/{channelID}/messages`  | `api.MessageChannel`  | Отправить новое сообщение в канал. |
|           | `/sse`                            | `api.CreateSSE`       | Создать Server-Sent Events (SSE) для взаимодействия с клиентами в режиме реального времени. |

Интерфейс:

```go
type APIInterface interface {
	GetAllChannels(w http.ResponseWriter, r *http.Request)
	AddNewChannel(w http.ResponseWriter, r *http.Request)
	DeleteChannel(w http.ResponseWriter, r *http.Request)
	GetChannel(w http.ResponseWriter, r *http.Request)
	GetMessagesChannel(w http.ResponseWriter, r *http.Request)
	MessageChannel(w http.ResponseWriter, r *http.Request)
	CreateSSE(w http.ResponseWriter, r *http.Request)
}
```

Регистрация обработчиков:

```go
r := mux.NewRouter()
r.HandleFunc("/channels", api.GetAllChannels).Methods("GET")
r.HandleFunc("/channels", api.AddNewChannel).Methods("POST")
r.HandleFunc("/channels/{channelID}", api.DeleteChannel).Methods("DELETE")
r.HandleFunc("/channels/{channelID}", api.GetChannel).Methods("GET")
r.HandleFunc("/channels/{channelID}/messages", api.GetMessagesChannel).Methods("GET")
r.HandleFunc("/channels/{channelID}/messages", api.MessageChannel).Methods("POST")
r.HandleFunc("/sse", api.CreateSSE)

http.Handle("/", r)
```
---

# База данных

В проекте для хранение данных использется MySQL.

Для использования необходимо: 

Установить [MySQL](https://dev.mysql.com/doc/refman/8.0/en/macos-installation.html)

Запустить сервер:
```bash
 mysql.server start
 # mysql.server stop
```

Подключиться к серверу: (root - иня пользователя mysql, он должен быть предварительно создан)
```bash
mysql -u root -p -h 127.0.0.1
```
В командной строке подключится к базе данный в даннос случае channel (БД должна быть предварительно создана)
```mysql
mysql> use channel
```
Запистить скрипт по созданию таблиц.
Файл **create-channels.sql** находится в папке `/MixChatServer/server/sql_db` и при запуске необходимо указать путь до него

```mysql
mysql> source create-channels.sql
# Query OK, 0 rows affected (0.01 sec)
# Query OK, 0 rows affected (0.01 sec)
# Query OK, 0 rows affected (0.00 sec)
# Query OK, 0 rows affected (0.01 sec)
```
Проверить, что таблица создана
```mysql
mysql> select * from channel;
# Empty set (0.01 sec)

mysql> select * from message;
# Empty set (0.00 sec)
```

Файл описывающий создание таблиц:

<pre>
-- Удаление таблицы message
DROP TABLE IF EXISTS message;

-- Удаление таблицы channel
DROP TABLE IF EXISTS channel;

CREATE TABLE channel (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logoURL VARCHAR(255),
    lastMessage VARCHAR(255),
    lastActivity DATETIME
);

CREATE TABLE message (
    id VARCHAR(255) PRIMARY KEY,
    text TEXT NOT NULL,
    userID VARCHAR(255) NOT NULL,
    userName VARCHAR(255) NOT NULL,
    date DATETIME NOT NULL,
    channelID VARCHAR(255) NOT NULL,
    FOREIGN KEY (channelID) REFERENCES channel(id) ON DELETE CASCADE ON UPDATE CASCADE
);
</pre>
