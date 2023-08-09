CREATE TABLE channel (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logoURL VARCHAR(255),
    lastMessage VARCHAR(255),
    lastActivity TIMESTAMP
);

CREATE TABLE message (
    id VARCHAR(255) PRIMARY KEY,
    text TEXT,
    userID VARCHAR(255) NOT NULL,
    userName VARCHAR(255) NOT NULL,
    date TIMESTAMP,
    channelID VARCHAR(255),
    FOREIGN KEY (channelID) REFERENCES channel(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Добавление стартовых данных в таблицу channel
INSERT INTO channel (id, name, logoURL, lastMessage, lastActivity)
VALUES
    ('channel1', 'Channel 1', 'logo1.jpg', 'Hello!', '2023-08-01 10:00:00'),
    ('channel2', 'Channel 2', 'logo2.jpg', 'Welcome!', '2023-08-02 12:30:00');

-- Добавление стартовых данных в таблицу message
INSERT INTO message (id, text, userID, userName, date, channelID)
VALUES
    ('message1', 'First message in Channel 1', 'user1', 'User A', '2023-08-01 10:15:00', 'channel1'),
    ('message2', 'Second message in Channel 1', 'user2', 'User B', '2023-08-01 11:30:00', 'channel1'),
    ('message3', 'Welcome message in Channel 2', 'user3', 'User C', '2023-08-02 12:45:00', 'channel2'),
    ('message4', 'Reply in Channel 2', 'user4', 'User D', '2023-08-02 13:00:00', 'channel2');

