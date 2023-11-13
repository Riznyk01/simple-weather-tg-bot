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
