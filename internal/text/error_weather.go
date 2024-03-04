package text

const (
	MoreInfoURLFormat      = "\n\n<a href=\"https://openweathermap.org/city/%s\">🌐 More</a>"
	FailedToGetWeather     = "Failed to get weather data:"
	TryAnother             = "Please try another city name, or try sending the location."
	ErrWhileGettingWeather = "Error occurred while getting weather JSON data."
	ErrWhileGeneratingURL  = "Error occurred while generating weather url."
	ErrDecodingJSON        = "Error occurred while decoding weather JSON data."
	ErrParsingWeatherURL   = "Can't parse the weather API URL."
	ErrFetchingWeather     = "Can't fetch weather by this city name."
)
