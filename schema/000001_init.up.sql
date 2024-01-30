CREATE TABLE IF NOT EXISTS user_data
(
    id BIGINT PRIMARY KEY,
    city VARCHAR(255),
    lat VARCHAR(20),
    lon VARCHAR(20),
    metric BOOLEAN,
    last VARCHAR(255)
);