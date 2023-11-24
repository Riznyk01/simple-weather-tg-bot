package utils

import (
	"math"
	"time"
)

func PressureConverting(hPa float64, metric bool) float64 {
	if metric {
		//HPaToMmHg
		return hPa * 0.750061561303
	}
	//HPaToInHg
	return hPa * 0.0295299830714
}
func TimeStampToHuman(timeStamp, timezone int, format string) string {
	location := time.FixedZone("Custom Timezone", timezone)
	timeValue := time.Unix(int64(timeStamp), 0).In(location)
	return timeValue.Format(format)
}
func TemperatureConverting(temp float64, metric bool) int16 {
	if metric {
		return int16(math.Round(temp))
	}
	return int16(math.Round((temp-273.15)*(9/5) + 32))
}
func TimeStampToInfo(timeStamp, timezone int, infoType string) string {
	location := time.FixedZone("Custom Timezone", timezone)
	timeValue := time.Unix(int64(timeStamp), 0).In(location)

	switch infoType {
	case "d":
		return timeValue.Weekday().String()
	case "m":
		return timeValue.Month().String()
	default:
		return "Invalid info type"
	}
}
func ToMilesPerHour(metersPerSecond float64) float64 {
	return metersPerSecond * 2.23694
}
func DegreesToDirectionIcon(degrees float64) string {
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
