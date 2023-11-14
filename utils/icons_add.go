package utils

func ReplaceWeatherToIcons(weather string) string {
	switch weather {
	case "scattered clouds":
		return "☁️"
	case "light rain":
		return "🌧️"
	case "moderate rain":
		return "🌧️"
	case "heavy intensity rain":
		return "🌧️🌧️"
	case "very heavy rain":
		return "🌧️🌧️🌧️"
	case "overcast clouds":
		return "🌥️"
	case "few clouds":
		return "☁️"
	case "broken clouds":
		return "🌦️"
	case "light snow":
		return "🌨️"
	case "clear sky":
		return "☀️"
	case "snow":
		return "❄️"
	default:
		return weather
	}
}
