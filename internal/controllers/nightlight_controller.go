package controllers

import (
	"luznocturna/luz-nocturna/internal/models"
	"luznocturna/luz-nocturna/internal/system"
)

// NightLightController maneja la lógica de la aplicación de luz nocturna
type NightLightController struct {
	config       *models.NightLightConfig
	appConfig    *models.AppConfig
	gammaManager *system.GammaManager
}

// NewNightLightController crea un nuevo controlador
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
