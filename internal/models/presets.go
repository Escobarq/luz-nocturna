package models

// TemperaturePresets define presets comunes de temperatura
type TemperaturePresets struct{}

var Presets = TemperaturePresets{}

// Presets de temperatura comunes
const (
	// Temperaturas predefinidas
	CandleLightTemp  = 3000 // Luz de vela
	WarmWhiteTemp    = 3500 // Blanco cálido
	NeutralWhiteTemp = 4500 // Blanco neutro
	CoolWhiteTemp    = 5500 // Blanco frío
	DaylightTemp     = 6500 // Luz diurna
)

// GetPresetName devuelve el nombre del preset más cercano a la temperatura dada
func (p TemperaturePresets) GetPresetName(temp float64) string {
	switch {
	case temp <= 3200:
		return "Muy cálida (🕯️)"
	case temp <= 3800:
		return "Cálida (🌅)"
	case temp <= 4800:
		return "Neutra (☀️)"
	case temp <= 6000:
		return "Fría (🌤️)"
	default:
		return "Diurna (☀️)"
	}
}

// GetRecommendedForTime devuelve una temperatura recomendada basada en la hora
func (p TemperaturePresets) GetRecommendedForTime(hour int) float64 {
	switch {
	case hour >= 22 || hour <= 6: // Noche
		return CandleLightTemp
	case hour >= 7 && hour <= 9: // Mañana
		return WarmWhiteTemp
	case hour >= 10 && hour <= 16: // Día
		return DaylightTemp
	case hour >= 17 && hour <= 21: // Tarde/Noche
		return NeutralWhiteTemp
	default:
		return NeutralWhiteTemp
	}
}
