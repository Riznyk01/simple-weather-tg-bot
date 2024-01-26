package util

import "fmt"

func WeatherTextToIcon(weather string, text bool) string {
	var weatherText string
	if text {
		weatherText = fmt.Sprintf(" [%s]", weather)
	}
	switch weather {
	case "scattered clouds":
		return "☁️" + weatherText
	case "light rain":
		return "🌧" + weatherText
	case "moderate rain":
		return "🌧" + weatherText
	case "heavy intensity rain":
		return "🌧🌧" + weatherText
	case "very heavy rain":
		return "🌧🌧🌧" + weatherText
	case "overcast clouds":
		return "🌥" + weatherText
	case "few clouds":
		return "☁️" + weatherText
	case "broken clouds":
		return "🌦" + weatherText
	case "light snow":
		return "🌨" + weatherText
	case "clear sky":
		return "☀️" + weatherText
	case "snow":
		return "❄️" + weatherText
	default:
		return weather
	}
}
