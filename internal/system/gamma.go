package system

import (
	"fmt"
	"os/exec"
)

// GammaManager maneja la configuración de gamma del sistema
type GammaManager struct{}

// NewGammaManager crea un nuevo manejador de gamma
func NewGammaManager() *GammaManager {
	return &GammaManager{}
}

// ApplyTemperature aplica una temperatura de color al sistema
func (gm *GammaManager) ApplyTemperature(temperature float64) error {
	// Verificar si redshift está disponible
	if err := gm.checkRedshiftAvailable(); err != nil {
		return fmt.Errorf("redshift no disponible: %w", err)
	}

	// Comando para aplicar temperatura con redshift
	// redshift -O <temperatura> -g <gamma>
	cmd := exec.Command("redshift", "-O", fmt.Sprintf("%.0f", temperature))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al aplicar temperatura %v: %w", temperature, err)
	}

	fmt.Printf("Temperatura aplicada: %.0fK\n", temperature)
	return nil
}

// Reset resetea la configuración de gamma a valores normales
func (gm *GammaManager) Reset() error {
	// Resetear con redshift
	cmd := exec.Command("redshift", "-x")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al resetear gamma: %w", err)
	}

	fmt.Println("Gamma reseteada a valores normales")
	return nil
}

// checkRedshiftAvailable verifica si redshift está instalado en el sistema
func (gm *GammaManager) checkRedshiftAvailable() error {
	cmd := exec.Command("which", "redshift")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("redshift no está instalado. Instálalo con: sudo apt install redshift")
	}
	return nil
}

// IsRedshiftInstalled verifica si redshift está instalado
func (gm *GammaManager) IsRedshiftInstalled() bool {
	return gm.checkRedshiftAvailable() == nil
}

// Alternative methods for other gamma tools

// ApplyWithXrandr aplica gamma usando xrandr (método alternativo)
func (gm *GammaManager) ApplyWithXrandr(temperature float64) error {
	// Convertir temperatura a valores RGB aproximados
	r, g, b := gm.temperatureToRGB(temperature)

	// Obtener displays disponibles
	displays, err := gm.getDisplays()
	if err != nil {
		return err
	}

	// Aplicar gamma a cada display
	for _, display := range displays {
		cmd := exec.Command("xrandr", "--output", display, "--gamma", fmt.Sprintf("%.2f:%.2f:%.2f", r, g, b))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error aplicando gamma a %s: %w", display, err)
		}
	}

	return nil
}

// getDisplays obtiene la lista de displays conectados
func (gm *GammaManager) getDisplays() ([]string, error) {
	cmd := exec.Command("xrandr", "--listactivemonitors")
	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Aquí iría el parsing del output de xrandr
	// Por simplicidad, devolvemos un display por defecto
	return []string{"eDP-1"}, nil // Display común en laptops
}

// temperatureToRGB convierte temperatura Kelvin a valores RGB gamma
func (gm *GammaManager) temperatureToRGB(temp float64) (r, g, b float64) {
	// Algoritmo simplificado para convertir temperatura a RGB
	// Valores entre 0.1 y 1.0

	if temp <= 3000 {
		r, g, b = 1.0, 0.4, 0.2
	} else if temp <= 4000 {
		r, g, b = 1.0, 0.7, 0.4
	} else if temp <= 5000 {
		r, g, b = 1.0, 0.9, 0.7
	} else {
		r, g, b = 1.0, 1.0, 1.0
	}

	return r, g, b
}
