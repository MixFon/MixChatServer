-- Удаление таблицы message
DROP TABLE IF EXISTS message;

-- Удаление таблицы channel
DROP TABLE IF EXISTS channel;

CREATE TABLE channel (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logoURL VARCHAR(255),
    lastMessage VARCHAR(255),
    lastActivity TIMESTAMP
);

CREATE TABLE message (
    id VARCHAR(255) PRIMARY KEY,
    text TEXT NOT NULL,
    userID VARCHAR(255) NOT NULL,
    userName VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    channelID VARCHAR(255) NOT NULL,
    FOREIGN KEY (channelID) REFERENCES channel(id) ON DELETE CASCADE ON UPDATE CASCADE
);


