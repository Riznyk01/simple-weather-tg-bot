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

func ReplaceWeatherPlusIcons(weather string) string {
	switch weather {
	case "scattered clouds":
		return "☁️ [scattered clouds]"
	case "light rain":
		return "🌧️ [rain]"
	case "moderate rain":
		return "🌧️ [moderate rain]"
	case "heavy intensity rain":
		return "🌧️🌧️ [heavy intensity rain]"
	case "very heavy rain":
		return "🌧️🌧️🌧️ [very heavy rain]"
	case "overcast clouds":
		return "🌥️ [overcast clouds]"
	case "few clouds":
		return "☁️ [few clouds]"
	case "broken clouds":
		return "🌦️ [broken clouds]"
	case "light snow":
		return "🌨️ [light snow]"
	case "clear sky":
		return "☀️ [clear sky]"
	case "snow":
		return "❄️ [snow]"
	default:
		return weather
	}
}
