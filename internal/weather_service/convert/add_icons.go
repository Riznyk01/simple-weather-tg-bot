package convert

import "fmt"

// AddIcon adds an icon corresponding to the weather condition.
// If 'text' is true, it appends the weather description in square brackets.
func AddIcon(weather string, text bool) string {
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

// DegsToDirIcon add icon with wind direction.
func DegsToDirIcon(degrees float64) string {
	if degrees >= 337.5 || degrees < 22.5 {
		return "⬆️"
	} else if degrees >= 22.5 && degrees < 67.5 {
		return "↗️"
	} else if degrees >= 67.5 && degrees < 112.5 {
		return "➡️"
	} else if degrees >= 112.5 && degrees < 157.5 {
		return "↘️"
	} else if degrees >= 157.5 && degrees < 202.5 {
		return "⬇️"
	} else if degrees >= 202.5 && degrees < 247.5 {
		return "↙️"
	} else if degrees >= 247.5 && degrees < 292.5 {
		return "⬅️"
	} else if degrees >= 292.5 && degrees < 337.5 {
		return "↖️"
	}
	return "Cannot determine"
}
