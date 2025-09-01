package models

import (
	"fmt"
)

// NightLightConfig representa la configuración de luz nocturna
type NightLightConfig struct {
	Temperature float64 // Temperatura en Kelvin
	MinTemp     float64 // Temperatura mínima
	MaxTemp     float64 // Temperatura máxima
	IsActive    bool    // Si está activa la luz nocturna
}

// NewNightLightConfig crea una nueva configuración con valores por defecto
func NewNightLightConfig() *NightLightConfig {
	return &NightLightConfig{
		Temperature: 4500, // Valor por defecto
		MinTemp:     3000, // Temperatura más cálida
		MaxTemp:     6500, // Temperatura más fría (luz diurna)
		IsActive:    false,
	}
}

// SetTemperature establece la temperatura asegurándose de que esté en el rango válido
func (config *NightLightConfig) SetTemperature(temp float64) {
	if temp < config.MinTemp {
		config.Temperature = config.MinTemp
	} else if temp > config.MaxTemp {
		config.Temperature = config.MaxTemp
	} else {
		config.Temperature = temp
	}
}

// GetTemperatureString devuelve la temperatura como string con formato
func (config *NightLightConfig) GetTemperatureString() string {
	return fmt.Sprintf("%.0fK", config.Temperature)
}

// Reset restablece la configuración a valores por defecto
func (config *NightLightConfig) Reset() {
	config.Temperature = 6500 // Luz diurna normal
	config.IsActive = false
}

// Apply activa la configuración de luz nocturna
func (config *NightLightConfig) Apply() error {
	config.IsActive = true
	// Aquí iría la lógica para aplicar realmente el filtro gamma
	// Por ahora solo marcamos como activa
	fmt.Printf("Aplicando luz nocturna con temperatura: %s\n", config.GetTemperatureString())
	return nil
}

// Disable desactiva la luz nocturna
func (config *NightLightConfig) Disable() error {
	config.IsActive = false
	// Aquí iría la lógica para desactivar el filtro gamma
	fmt.Println("Desactivando luz nocturna")
	return nil
}
