CREATE TABLE IF NOT EXISTS schedule_data
(
    id BIGINT,
    city VARCHAR(255),
    schedule_time TIMESTAMP,
    weather_type VARCHAR(10),
    timezone_offset FLOAT
);