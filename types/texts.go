package types

// Constants for commands
const (
	CommandStart            = "/start"
	CommandHelp             = "/help"
	CommandCurrent          = "current"
	CommandForecast         = "5-days forecast"
	CommandForecastLocation = "5-days forecast 📍"
	CommandCurrentLocation  = "current 📍"
)

// Constants for messages
const (
	WelcomeMessage            = "Hello! This bot will send you weather information from openweathermap.org. "
	HelpMessage               = "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
	MissingCityMessage        = "You didn't enter a city.\nPlease enter a city or send your location,\nand then choose the type of weather."
	ChooseOptionMessage       = "Choose an action:"
	NoLocationProvidedMessage = "You tried to get the weather based on your location, but you didn't share your location."
)
