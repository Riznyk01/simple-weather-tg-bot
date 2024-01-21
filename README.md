# Simple weather tg bot

Simple weather tg bot is a Telegram bot written in Golang that provides weather information from the free OpenWeatherMap API. Users can input city names or send their locations to receive current weather details on demand.

## Table of Contents

- [Environment Variables](#environment-variables)
- [Installation](#installation)
- [Usage](#usage)
- [Bot Commands](#bot-commands)
- [Pricing Information](#pricing-information)
- [Note](#note)
- [Screenshots](#Screenshots)

## Screenshots

Here are some screenshots of the bot in action:

*City name current weather search:*  
![Screenshot 1](images/screenshot_current.png)

*City name 5-days forecast search:*  
![Screenshot 2](images/screenshot_5day_forecast.png)

*Weather retrieval based on location:*  
![Screenshot 3](images/screenshot_location.png)

## Environment Variables

Before running the bot, make sure to set the following environment variables:

BOT_DEBUG=(true or false): "Enable/disable debug mode of the Telegram Bot API."  
BOT_TOKEN=YOUR_BOT_TOKEN: "Your Telegram Bot Token."  
LOG_LEVEL=(panic/fatal/error/warn/warning/info/debug/trace): "Sets the log level."  
TYPE_OF_LOG=(JSONLOG or TEXTLOG): "Choose between JSON or text log format."    
WEATHER_API_URL=https://api.openweathermap.org/data/2.5/: "OpenWeatherMap API URL."  
WEATHER_TOKEN=YOUR_OPENWEATHERMAP_API_KEY: "Your OpenWeatherMap API Key."  

Replace `YOUR_BOT_TOKEN` with your Telegram Bot Token, which you can obtain by creating a new bot on Telegram. Follow these steps:

1. Open Telegram and search for the "BotFather" bot (@BotFather).
2. Start a chat with BotFather and use the `/newbot` command to create a new bot.
3. Follow the instructions from BotFather to choose a name and username for your bot.
4. Once the bot is created, BotFather will provide you with a token.

Get your free OpenWeatherMap API Key [here](https://home.openweathermap.org/api_keys).

Set the values of the LOG_LEVEL and TYPE_OF_LOG variables from the provided options.

## Installation

Clone the repository:

```bash
git clone https://github.com/Riznyk01/simple-weather-tg-bot.git
cd simple-weather-tg-bot
```

Build the Docker image:
```bash
sudo docker build -t simple-weather-tg-bot .
```
## Usage
Run the Docker container:
```bash
sudo docker run -p 8080:8080 \
    -e BOT_DEBUG=... \
    -e BOT_TOKEN=... \
    -e LOG_LEVEL=... \
    -e TYPE_OF_LOG=... \
    -e WEATHER_API_URL=https://api.openweathermap.org/data/2.5/ \
    -e WEATHER_TOKEN=... \
    simple-weather-tg-bot:latest
```


## Bot Commands
/start: sends a welcome message and instructions to the user.  
/help: provides information on how to use the bot.  
/metric: set metric units  
/nonmetric: set non-metric units

## Pricing Information
Information about pricing, available plans, and limitations of the free API package can be found [here](https://openweathermap.org/price).

## Note
Some cities may return weather information correctly using the city name, while others may require the user's location. Use the preferred option for accurate results.

Feel free to contribute and enhance the functionality of this bot!