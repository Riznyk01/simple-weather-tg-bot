# Simple weather tg bot

Simple weather tg bot is a Telegram bot written in Golang that provides weather information from the free OpenWeatherMap API. Users can input city names or send their locations to receive current weather details on demand.

## Table of Contents

- [Environment Variables](#environment-variables)
- [Installation](#installation)
- [Usage](#usage)
- [Bot Commands](#bot-commands)
- [Pricing Information](#pricing-information)
- [Note](#note)

## Environment Variables

Make sure to set the following environment variables before running the bot. Create a `.env.dev` file in the root directory with the following content:

BOT_TOKEN=YOUR_BOT_TOKEN

WEATHER_KEY=YOUR_OPENWEATHERMAP_API_KEY

Replace `YOUR_BOT_TOKEN` with your Telegram Bot Token, which you can obtain by creating a new bot on Telegram. Follow these steps:

1. Open Telegram and search for the "BotFather" bot (@BotFather).
2. Start a chat with BotFather and use the `/newbot` command to create a new bot.
3. Follow the instructions from BotFather to choose a name and username for your bot.
4. Once the bot is created, BotFather will provide you with a token. Copy the token and replace `YOUR_BOT_TOKEN` in the `.env.dev` file.

Get your free OpenWeatherMap API Key [here](https://home.openweathermap.org/api_keys).

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
sudo docker run -d simple-weather-tg-bot
```

Check the running containers:

```bash
sudo docker ps
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

Feel free to contribute and enhance the functionality of this simple weather bot!