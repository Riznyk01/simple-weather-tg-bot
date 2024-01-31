package model

const (
	MessageWelcome             = "🎈Hello, %s. This bot will send you weather information from openweathermap.org. \n\n"
	MessageHelp                = "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
	MessageChooseOption        = "Choose an action:"
	MessageMetricUnitChanged   = "The unit system has been updated."
	MessageLastDataUnavailable = "Sorry ❤️, %s, there is no saved weather forecast parameters from your last request. <b>" +
		"\n\nPlease try sending the city name or location, and then select the desired weather type using the buttons.</b>"
	MessageSetUsersSystemError    = "Error while saving user's preferred system of measurement."
	MessageSetUsersLocationError  = "Error while saving user's preferred location."
	MessageSetUsersCityError      = "Error while saving user's preferred city."
	MessageUnsupportedMessageType = "Sorry, this type of message is not supported by the bot."
	MessageUnsupportedCommand     = "Sorry, this command is not supported by the bot."
)
