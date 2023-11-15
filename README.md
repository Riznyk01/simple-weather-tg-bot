# SimpleWeatherTgBot

SimpleWeatherTgBot is a Telegram bot written in Golang that provides weather information from the free OpenWeatherMap API. Users can input city names or send their locations to get weather updates.

## Table of Contents

- [Environment Variables](#environment-variables)
- [Installation](#installation)
- [Usage](#usage)
- [Bot Commands](#bot-commands)
- [Note](#note)

## Environment Variables

Make sure to set the following environment variables before running the bot. Create a `.env.dev` file in the root directory with the following content:

BOT_TOKEN=YOUR_BOT_TOKEN

WEATHER_KEY=YOUR_OPENWEATHERMAP_API_KEY

Replace `YOUR_BOT_TOKEN` with your Telegram Bot Token and `YOUR_OPENWEATHERMAP_API_KEY` with your free OpenWeatherMap API Key.

## Installation

Clone the repository:

```bash
git clone https://github.com/Riznyk01/SimpleWeatherTgBot.git
cd SimpleWeatherTgBot
```

Build the Docker image:
```bash
sudo docker build -t simpleweathertgbot .
```
## Usage
Run the Docker container:
```bash
sudo docker run -d simpleweathertgbot
```

Check the running containers:

```bash
sudo docker ps
```
## Bot Commands
/start: Sends a welcome message and instructions to the user.
/help: Provides information on how to use the bot.
## Note
Some cities may return weather information correctly using the city name, while others may require the user's location. Use the preferred option for accurate results.

Feel free to contribute and enhance the functionality of this simple weather bot!