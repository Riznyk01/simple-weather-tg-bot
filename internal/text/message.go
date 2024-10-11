package text

const (
	MsgWelcome             = "🎈Hello, %s. This bot will send you weather information from openweathermap.org. \n\n"
	MsgHelp                = "Enter the city name in your language, then choose the weather type, or send your location, and then also choose the weather type. \n\nBot Commands for Adding Schedules:\n🪄The command \"/add_18:00_2_cityname_weathertype_metricunits\" adds a schedule that will be executed at 18:00, when:\n\n🔸 2 — the user's timezone\n🔸 weathertype — the forecast type (\"current\", \"5-days forecast\", \"today forecast\"),\n🔸 metricunits — true/false\n\n🪄The command \"/deleteschedules\" deletes all user schedules.\n\n🪄The command \"/viewschedules\" fetches all user schedules."
	MsgChooseOption        = "Choose an action:"
	MsgMetricUnitChanged   = "The unit system has been updated."
	MsgLastDataUnavailable = "Sorry ❤️, %s, there is no saved weather forecast parameters from your last request. <b>" +
		"\n\nPlease try sending the city name or location, and then select the desired weather type using the buttons.</b>"
	MsgSetUsersSystemError    = "Error while saving user's preferred system of measurement."
	MsgSetUsersLocationError  = "Error while saving user's preferred location."
	MsgSetUsersCityError      = "Error while saving user's preferred city."
	MsgUnsupportedMessageType = "Sorry, this type of message is not supported by the bot."
	MsgAlreadyStarted         = "User already started the bot."
	MsgUnsupportedCommand     = "Sorry, this command is not supported by the bot."
)
