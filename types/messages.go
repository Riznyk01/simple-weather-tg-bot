package types

const (
	WelcomeMessage      = "🎈Hello, %s. This bot will send you weather information from openweathermap.org. \n\n"
	HelpMessage         = "Enter the city name in any language, then choose the weather type, or send your location, and then also choose the weather type."
	ChooseOptionMessage = "Choose an action:"
	MetricUnitOn        = "Metric units are enabled."
	MetricUnitOff       = "Metric units are disabled."
	LastDataUnavailable = "Sorry ❤️, %s, the forecast with the latest parameters is unavailable due to a bot restart. <b>" +
		"\n\nPlease try sending the city name or location, and then select the weather type using the buttons.</b>"
	SetUsersSystemError   = "Error while saving user's preferred system of measurement."
	SetUsersLocationError = "Error while saving user's preferred location."
	SetUsersCityError     = "Error while saving user's preferred city."
)
