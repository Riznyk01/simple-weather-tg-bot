package convert

import (
	"math"
)

// Pressure converts the pressure [hPa] to [InHg] or [mmHg].
func Pressure(hPa int, metric bool) int {
	if metric {
		return int(float64(hPa) * 0.750061561303) //[mmHg]
	}
	return int(float64(hPa) * 0.0295299830714) //[inHg]
}

// KelvinToFahrenheitAndRound converts temperature from Kelvin to Fahrenheit and rounds it to the nearest integer.
func KelvinToFahrenheitAndRound(t float64, metric bool) int16 {
	if metric {
		return int16(math.Round(t))
	}
	return int16(math.Round((t-273.15)*(9/5) + 32))
}

// WindSpeed returns wind speed in [m/s] or [mph] depending on metric/non-metric units.
func WindSpeed(ms float64, metric bool) float64 {
	if metric {
		return ms
	}
	return ms * 2.23694
}

// Units returns units based on the metric system.
func Units(metricUnits bool) (tempUnits, windUnits, pressureUnits string) {
	if metricUnits {
		return "°C", "m/s", "mmHg"
	}
	return "°F", "mph", "inHg"
}
