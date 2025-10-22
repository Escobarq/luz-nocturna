package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppConfig representa la configuración persistente de la aplicación
type AppConfig struct {
	LastTemperature float64        `json:"last_temperature"`
	AutoStart       bool           `json:"auto_start"`
	MinimizeToTray  bool           `json:"minimize_to_tray"`
	StartMinimized  bool           `json:"start_minimized"`
	ScheduleEnabled bool           `json:"schedule_enabled"`
	Schedule        ScheduleConfig `json:"schedule"`
}

// ScheduleConfig representa la configuración de horarios automáticos
type ScheduleConfig struct {
	StartTime          string  `json:"start_time"`           // Formato "HH:MM" para inicio del filtro nocturno
	EndTime            string  `json:"end_time"`             // Formato "HH:MM" para fin del filtro nocturno
	NightTemp          float64 `json:"night_temp"`           // Temperatura nocturna (ej: 3000K)
	DayTemp            float64 `json:"day_temp"`             // Temperatura diurna (ej: 6500K)
	TransitionTime     int     `json:"transition_time"`      // Tiempo de transición en minutos
	AutoDetectLocation bool    `json:"auto_detect_location"` // Detectar ubicación para sunrise/sunset automático
}

// NewAppConfig crea una nueva configuración con valores por defecto
func NewAppConfig() *AppConfig {
	return &AppConfig{
		LastTemperature: 4500,
		AutoStart:       false,
		MinimizeToTray:  true,
		StartMinimized:  false,
		ScheduleEnabled: false,
		Schedule: ScheduleConfig{
			StartTime:          "20:00",
			EndTime:            "07:00",
			NightTemp:          3200,
			DayTemp:            6500,
			TransitionTime:     30,
			AutoDetectLocation: false,
		},
	}
}

// GetConfigPath devuelve la ruta del archivo de configuración
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "luz-nocturna", "config.json")
}

// Load carga la configuración desde el archivo
func (config *AppConfig) Load() error {
	configPath := GetConfigPath()

	// Crear directorio si no existe
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Si el archivo no existe, usar valores por defecto
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config.Save() // Crear archivo con valores por defecto
	}

	// Leer archivo
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Deserializar JSON
	return json.Unmarshal(data, config)
}

// Save guarda la configuración al archivo
func (config *AppConfig) Save() error {
	configPath := GetConfigPath()

	// Crear directorio si no existe
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Serializar a JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Escribir archivo
	return os.WriteFile(configPath, data, 0644)
}
