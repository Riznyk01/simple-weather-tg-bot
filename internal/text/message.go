package text

const (
	MsgWelcome             = "🎈Hello, %s. This bot will send you weather information from openweathermap.org. \n\n"
	MsgHelp                = "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
	MsgChooseOption        = "Choose an action:"
	MsgMetricUnitChanged   = "The unit system has been updated."
	MsgLastDataUnavailable = "Sorry ❤️, %s, there is no saved weather forecast parameters from your last request. <b>" +
		"\n\nPlease try sending the city name or location, and then select the desired weather type using the buttons.</b>"
	MsgSetUsersSystemError    = "Error while saving user's preferred system of measurement."
	MsgSetUsersLocationError  = "Error while saving user's preferred location."
	MsgSetUsersCityError      = "Error while saving user's preferred city."
	MsgUnsupportedMessageType = "Sorry, this type of message is not supported by the bot."
	MsgUnsupportedCommand     = "Sorry, this command is not supported by the bot."
)
