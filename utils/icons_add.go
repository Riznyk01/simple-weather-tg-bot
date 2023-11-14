package utils

func AddWeatherIcons(weather string) string {
	switch weather {
	case "Rain":
		return "🌧 Rain"
	case "Clouds":
		return "☁️ Clouds"
	case "Clear":
		return "✨ Clear"
	case "Snow":
		return "❄️ Snow"
	default:
		return weather
	}
}

func ReplaceWeatherToIcons(weather string) string {
	switch weather {
	case "scattered clouds":
		return "☁️"
	case "light rain":
		return "🌧️"
	case "moderate rain":
		return "🌧️"
	case "overcast clouds":
		return "🌥️"
	case "few clouds":
		return "☁️"
	case "broken clouds":
		return "🌦️"
	case "light snow":
		return "🌨️"
	default:
		return weather
	}
}
