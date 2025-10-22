package controllers

import (
	"fmt"
	"luznocturna/luz-nocturna/internal/models"
	"luznocturna/luz-nocturna/internal/system"
	"time"
)

/**
 * NightLightController - Controlador principal de la aplicación
 *
 * Maneja toda la lógica de negocio para el control de temperatura de color
 * del monitor. Coordina entre los modelos de configuración y el sistema
 * de gestión de gamma del display.
 *
 * @struct {NightLightController}
 * @property {*models.NightLightConfig} config - Configuración actual de luz nocturna
 * @property {*models.AppConfig} appConfig - Configuración persistente de la aplicación
 * @property {*system.GammaManager} gammaManager - Manejador de gamma del sistema
 */
type NightLightController struct {
	config       *models.NightLightConfig
	appConfig    *models.AppConfig
	gammaManager *system.GammaManager
	scheduler    *models.Scheduler
}

/**
 * NewNightLightController - Constructor del controlador principal
 *
 * Inicializa el controlador con configuración por defecto y carga
 * la configuración persistente si existe. Configura el manejador
 * de gamma del sistema según la plataforma detectada.
 *
 * @returns {*NightLightController} Nueva instancia del controlador
 *
 * @example
 *   controller := NewNightLightController()
 *   controller.ApplyNightLight()
 */
func NewNightLightController() *NightLightController {
	controller := &NightLightController{
		config:       models.NewNightLightConfig(),
		appConfig:    models.NewAppConfig(),
		gammaManager: system.NewGammaManager(),
	}

	// Cargar configuración guardada
	if err := controller.appConfig.Load(); err == nil {
		controller.config.SetTemperature(controller.appConfig.LastTemperature)
	}

	// Inicializar programador con callback para aplicar temperatura
	controller.scheduler = models.NewScheduler(controller.appConfig, func(temp float64) error {
		controller.config.SetTemperature(temp)
		return controller.gammaManager.ApplyTemperature(temp)
	})

	// Iniciar programación automática si está habilitada
	if controller.appConfig.ScheduleEnabled {
		controller.scheduler.Start()
	}

	return controller
}

// GetConfig devuelve la configuración actual
func (c *NightLightController) GetConfig() *models.NightLightConfig {
	return c.config
}

// GetAppConfig devuelve la configuración de la aplicación
func (c *NightLightController) GetAppConfig() *models.AppConfig {
	return c.appConfig
}

// UpdateTemperature actualiza la temperatura
func (c *NightLightController) UpdateTemperature(temp float64) {
	c.config.SetTemperature(temp)
	// Guardar la temperatura como preferencia del usuario
	c.appConfig.LastTemperature = temp
	c.appConfig.Save() // Ignorar errores por ahora
}

// ApplyNightLight aplica la configuración de luz nocturna usando xrandr
func (c *NightLightController) ApplyNightLight() error {
	// Aplicar temperatura usando nuestro sistema xrandr
	if err := c.gammaManager.ApplyTemperature(c.config.Temperature); err != nil {
		return err
	}

	// Marcar como aplicado en el modelo
	return c.config.Apply()
}

// ResetNightLight resetea la configuración a valores por defecto
func (c *NightLightController) ResetNightLight() error {
	// Resetear gamma del sistema
	if err := c.gammaManager.Reset(); err != nil {
		// Si falla, al menos resetear el modelo
		c.config.Reset()
		return err
	}

	// Resetear configuración
	c.config.Reset()
	c.appConfig.LastTemperature = c.config.Temperature
	c.appConfig.Save() // Ignorar errores

	return nil
}

// ToggleNightLight alterna entre activar y desactivar la luz nocturna
func (c *NightLightController) ToggleNightLight() error {
	if c.config.IsActive {
		return c.ResetNightLight()
	}
	return c.ApplyNightLight()
}

// GetTemperatureRange devuelve el rango de temperatura válido
func (c *NightLightController) GetTemperatureRange() (min, max float64) {
	return c.config.MinTemp, c.config.MaxTemp
}

// GetDisplays devuelve la lista de displays detectados
func (c *NightLightController) GetDisplays() []string {
	return c.gammaManager.GetDisplays()
}

// === MÉTODOS DE PROGRAMACIÓN AUTOMÁTICA ===

// EnableSchedule habilita la programación automática
func (c *NightLightController) EnableSchedule(enabled bool) {
	c.appConfig.ScheduleEnabled = enabled
	c.appConfig.Save()

	if enabled {
		c.scheduler.Start()
	} else {
		c.scheduler.Stop()
	}

	c.scheduler.UpdateConfig(c.appConfig)
}

// IsScheduleEnabled verifica si la programación está habilitada
func (c *NightLightController) IsScheduleEnabled() bool {
	return c.appConfig.ScheduleEnabled
}

// IsScheduleRunning verifica si el programador está ejecutándose
func (c *NightLightController) IsScheduleRunning() bool {
	return c.scheduler.IsRunning()
}

// UpdateScheduleConfig actualiza la configuración de horarios
func (c *NightLightController) UpdateScheduleConfig(startTime, endTime string, nightTemp, dayTemp float64, transitionTime int) {
	c.appConfig.Schedule.StartTime = startTime
	c.appConfig.Schedule.EndTime = endTime
	c.appConfig.Schedule.NightTemp = nightTemp
	c.appConfig.Schedule.DayTemp = dayTemp
	c.appConfig.Schedule.TransitionTime = transitionTime
	c.appConfig.Save()

	c.scheduler.UpdateConfig(c.appConfig)
}

// GetScheduleConfig obtiene la configuración actual de horarios
func (c *NightLightController) GetScheduleConfig() models.ScheduleConfig {
	return c.appConfig.Schedule
}

// GetNextScheduleChange obtiene información sobre el próximo cambio programado
func (c *NightLightController) GetNextScheduleChange() (string, float64, time.Duration) {
	return c.scheduler.GetNextScheduleChange()
}

// ApplyScheduleNow aplica inmediatamente la temperatura correspondiente al horario actual
func (c *NightLightController) ApplyScheduleNow() error {
	if !c.appConfig.ScheduleEnabled {
		return fmt.Errorf("la programación automática está deshabilitada")
	}

	// El scheduler aplicará automáticamente la temperatura correcta
	c.scheduler.Stop()
	c.scheduler.Start()
	return nil
}
