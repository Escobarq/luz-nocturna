package system

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

/**
 * GammaManager - Manejador principal del sistema de gamma
 * 
 * Maneja la configuración de temperatura de color del sistema
 * soportando tanto X11 (xrandr) como Wayland (wlr-gamma-control).
 * 
 * @struct {GammaManager}
 * @property {[]string} displays - Lista de displays detectados automáticamente
 * @property {string} protocol - Protocolo de display detectado ("x11" o "wayland")
 */
type GammaManager struct {
	displays []string
	protocol string
}

/**
 * NewGammaManager - Constructor del manejador de gamma
 * 
 * Inicializa un nuevo manejador de gamma, detecta automáticamente
 * el protocolo de display (X11/Wayland) y los displays disponibles.
 * 
 * @returns {*GammaManager} Nueva instancia del manejador de gamma
 * @example
 *   gm := NewGammaManager()
 *   gm.ApplyTemperature(4000) // Aplica 4000K
 */
func NewGammaManager() *GammaManager {
	gm := &GammaManager{}
	gm.detectDisplayProtocol()
	gm.detectDisplays()
	return gm
}

/**
 * ApplyTemperature - Aplica una temperatura de color específica
 * 
 * Convierte la temperatura en Kelvin a valores RGB gamma y los aplica
 * a todos los displays detectados usando el protocolo apropiado.
 * 
 * @param {float64} temperature - Temperatura en Kelvin (3000-6500)
 * @returns {error} Error si no se puede aplicar la temperatura
 * @example
 *   err := gm.ApplyTemperature(3500) // Temperatura cálida
 *   if err != nil {
 *       log.Printf("Error: %v", err)
 *   }
 */
func (gm *GammaManager) ApplyTemperature(temperature float64) error {
	// Convertir temperatura a valores RGB gamma
	r, g, b := gm.temperatureToRGB(temperature)
	
	if gm.protocol == "wayland" {
		return gm.applyWaylandGamma(r, g, b)
	}
	
	// Aplicar usando X11/xrandr (comportamiento por defecto)
	return gm.applyX11Gamma(r, g, b, temperature)
}

/**
 * Reset - Resetea la configuración de gamma a valores normales
 * 
 * Restaura todos los displays a gamma normal (1.0:1.0:1.0),
 * removiendo cualquier filtro de temperatura de color aplicado.
 * 
 * @returns {error} Error si no se puede resetear
 * @example
 *   err := gm.Reset()
 *   if err != nil {
 *       log.Printf("No se pudo resetear: %v", err)
 *   }
 */
func (gm *GammaManager) Reset() error {
	if gm.protocol == "wayland" {
		return gm.resetWaylandGamma()
	}
	
	// Reset usando X11/xrandr
	for _, display := range gm.displays {
		cmd := exec.Command("xrandr", "--output", display, "--gamma", "1.0:1.0:1.0")
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Advertencia: no se pudo resetear gamma en %s: %v\n", display, err)
			continue
		}
	}
	
	fmt.Println("✅ Gamma reseteada a valores normales")
	return nil
}

/**
 * detectDisplayProtocol - Detecta el protocolo de display en uso
 * 
 * Determina si el sistema está ejecutando X11 o Wayland
 * verificando variables de entorno y procesos activos.
 * 
 * @private
 */
func (gm *GammaManager) detectDisplayProtocol() {
	// Verificar variables de entorno
	if os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland" {
		gm.protocol = "wayland"
		return
	}
	
	// Por defecto asumir X11
	gm.protocol = "x11"
}

/**
 * detectDisplays - Detecta automáticamente los displays conectados
 * 
 * Escanea el sistema para encontrar todos los displays/monitores
 * conectados usando las herramientas apropiadas según el protocolo.
 * 
 * @private
 */
func (gm *GammaManager) detectDisplays() {
	if gm.protocol == "wayland" {
		gm.detectWaylandDisplays()
		return
	}
	
	// Detectar displays X11 usando xrandr
	cmd := exec.Command("xrandr")
	output, err := cmd.Output()
	if err != nil {
		// Fallback a display común
		gm.displays = []string{"eDP-1"}
		fmt.Printf("⚠️  No se pudo ejecutar xrandr, usando display por defecto: eDP-1\n")
		return
	}
	
	// Parsear output de xrandr para encontrar displays conectados
	lines := strings.Split(string(output), "\n")
	connectedRegex := regexp.MustCompile(`^(\S+)\s+connected`)
	
	var displays []string
	for _, line := range lines {
		if matches := connectedRegex.FindStringSubmatch(line); matches != nil {
			displays = append(displays, matches[1])
		}
	}
	
	if len(displays) == 0 {
		// Fallback si no se detecta nada
		displays = []string{"eDP-1"}
	}
	
	gm.displays = displays
	fmt.Printf("🖥️  Displays detectados (%s): %v\n", gm.protocol, displays)
}

/**
 * applyX11Gamma - Aplica gamma usando xrandr (X11)
 * 
 * @param {float64} r - Componente rojo del gamma (0.3-1.0)
 * @param {float64} g - Componente verde del gamma (0.3-1.0) 
 * @param {float64} b - Componente azul del gamma (0.3-1.0)
 * @param {float64} temperature - Temperatura original para logging
 * @returns {error} Error si falla la aplicación
 * @private
 */
func (gm *GammaManager) applyX11Gamma(r, g, b, temperature float64) error {
	for _, display := range gm.displays {
		cmd := exec.Command("xrandr", "--output", display, "--gamma", fmt.Sprintf("%.2f:%.2f:%.2f", r, g, b))
		if err := cmd.Run(); err != nil {
			// Si falla un display, continúa con los otros
			fmt.Printf("⚠️  Advertencia: no se pudo aplicar gamma a %s: %v\n", display, err)
			continue
		}
	}
	
	fmt.Printf("🌡️  Temperatura aplicada: %.0fK (RGB: %.2f:%.2f:%.2f)\n", temperature, r, g, b)
	return nil
}

/**
 * applyWaylandGamma - Aplica gamma usando wlr-gamma-control (Wayland)
 * 
 * Utiliza wl-gamma-relay o gammastep para aplicar temperatura de color
 * en entornos Wayland que soportan wlr-gamma-control-unstable-v1.
 * 
 * @param {float64} r - Componente rojo del gamma (0.3-1.0)
 * @param {float64} g - Componente verde del gamma (0.3-1.0)
 * @param {float64} b - Componente azul del gamma (0.3-1.0)
 * @returns {error} Error si falla la aplicación
 * @private
 */
func (gm *GammaManager) applyWaylandGamma(r, g, b float64) error {
	// Intentar con wl-gamma-relay primero
	cmd := exec.Command("wl-gamma-relay", fmt.Sprintf("%.2f", r), fmt.Sprintf("%.2f", g), fmt.Sprintf("%.2f", b))
	if err := cmd.Run(); err == nil {
		fmt.Printf("🌡️  Gamma aplicada en Wayland (wl-gamma-relay): %.2f:%.2f:%.2f\n", r, g, b)
		return nil
	}
	
	// Fallback: Intentar con wlsunset si está disponible
	cmd = exec.Command("pkill", "wlsunset")
	cmd.Run() // Matar instancia anterior si existe
	
	// Calcular temperatura aproximada desde RGB
	temp := gm.rgbToTemperature(r, g, b)
	cmd = exec.Command("wlsunset", "-t", fmt.Sprintf("%.0f", temp), "-T", fmt.Sprintf("%.0f", temp))
	if err := cmd.Start(); err == nil {
		fmt.Printf("🌡️  Temperatura aplicada en Wayland (wlsunset): %.0fK\n", temp)
		return nil
	}
	
	return fmt.Errorf("no se pudo aplicar gamma en Wayland - instala wl-gamma-relay o wlsunset")
}

/**
 * resetWaylandGamma - Resetea gamma en Wayland
 * 
 * @returns {error} Error si falla el reset
 * @private
 */
func (gm *GammaManager) resetWaylandGamma() error {
	// Matar procesos de control de gamma
	exec.Command("pkill", "wlsunset").Run()
	exec.Command("pkill", "wl-gamma-relay").Run()
	
	// Resetear con wl-gamma-relay
	cmd := exec.Command("wl-gamma-relay", "1.0", "1.0", "1.0")
	if err := cmd.Run(); err == nil {
		fmt.Println("✅ Gamma reseteada en Wayland")
		return nil
	}
	
	return fmt.Errorf("no se pudo resetear gamma en Wayland")
}

/**
 * detectWaylandDisplays - Detecta displays en Wayland
 * 
 * En Wayland, el control de gamma se aplica globalmente,
 * por lo que no necesitamos detectar displays específicos.
 * 
 * @private
 */
func (gm *GammaManager) detectWaylandDisplays() {
	// En Wayland, el control de gamma es global
	gm.displays = []string{"wayland-global"}
	fmt.Printf("🖥️  Protocolo Wayland detectado - control global de gamma\n")
}

/**
 * GetDisplays - Obtiene la lista de displays detectados
 * 
 * @returns {[]string} Lista de nombres de displays
 * @example
 *   displays := gm.GetDisplays()
 *   fmt.Printf("Displays disponibles: %v", displays)
 */
func (gm *GammaManager) GetDisplays() []string {
	return gm.displays
}

/**
 * GetProtocol - Obtiene el protocolo de display detectado
 * 
 * @returns {string} Protocolo detectado ("x11" o "wayland")
 */
func (gm *GammaManager) GetProtocol() string {
	return gm.protocol
}

/**
 * temperatureToRGB - Convierte temperatura Kelvin a valores RGB gamma
 * 
 * Implementa el algoritmo de Tanner Helland para conversión de temperatura
 * de color a valores RGB, optimizado para control de gamma en pantallas.
 * 
 * @param {float64} temp - Temperatura en Kelvin (1000-40000, típicamente 3000-6500)
 * @returns {float64, float64, float64} Componentes RGB normalizados (0.3-1.0)
 * @example
 *   r, g, b := gm.temperatureToRGB(4000) // Temperatura cálida
 *   // r ≈ 1.0, g ≈ 0.8, b ≈ 0.6
 */
func (gm *GammaManager) temperatureToRGB(temp float64) (r, g, b float64) {
	// Algoritmo de Tanner Helland optimizado para control de gamma
	// Basado en datos empíricos de temperatura de color de cuerpo negro
	
	// Normalizar temperatura (dividir por 100 para cálculos)
	temp = temp / 100

	// === CALCULAR COMPONENTE ROJO ===
	if temp <= 66 {
		// Para temperaturas <= 6600K, el rojo está al máximo
		r = 1.0
	} else {
		// Para temperaturas > 6600K, calcular curva de enfriamiento
		r = temp - 60
		r = 329.698727446 * math.Pow(r, -0.1332047592)
		if r < 0 {
			r = 0
		}
		if r > 1 {
			r = 1
		}
	}

	// === CALCULAR COMPONENTE VERDE ===
	if temp <= 66 {
		// Curva de calentamiento para verde
		g = temp
		g = 99.4708025861*math.Log(g) - 161.1195681661
		if g < 0 {
			g = 0
		}
		if g > 255 {
			g = 255
		}
		g = g / 255 // Normalizar a 0-1
	} else {
		// Curva de enfriamiento para verde
		g = temp - 60
		g = 288.1221695283 * math.Pow(g, -0.0755148492)
		if g < 0 {
			g = 0
		}
		if g > 1 {
			g = 1
		}
	}

	// === CALCULAR COMPONENTE AZUL ===
	if temp >= 66 {
		// Para temperaturas >= 6600K, el azul está al máximo
		b = 1.0
	} else if temp <= 19 {
		// Para temperaturas muy bajas, no hay azul
		b = 0
	} else {
		// Curva de calentamiento para azul
		b = temp - 10
		b = 138.5177312231*math.Log(b) - 305.0447927307
		if b < 0 {
			b = 0
		}
		if b > 255 {
			b = 255
		}
		b = b / 255 // Normalizar a 0-1
	}

	// === APLICAR LÍMITES MÍNIMOS PARA GAMMA ===
	// Evitar valores demasiado extremos que puedan dañar la vista
	// o hacer la pantalla ilegible
	const minGamma = 0.3
	if r < minGamma {
		r = minGamma
	}
	if g < minGamma {
		g = minGamma
	}
	if b < minGamma {
		b = minGamma
	}

	return r, g, b
}

/**
 * rgbToTemperature - Convierte valores RGB aproximadamente a temperatura Kelvin
 * 
 * Función inversa aproximada para estimar temperatura desde valores RGB.
 * Útil para retrocompatibilidad con herramientas que requieren temperatura.
 * 
 * @param {float64} r - Componente rojo (0-1)
 * @param {float64} g - Componente verde (0-1)
 * @param {float64} b - Componente azul (0-1)
 * @returns {float64} Temperatura estimada en Kelvin
 * @private
 */
func (gm *GammaManager) rgbToTemperature(r, g, b float64) float64 {
	// Estimación simple basada en la relación azul/rojo
	ratio := b / r
	
	if ratio >= 0.9 {
		return 6500 // Temperatura diurna
	} else if ratio >= 0.7 {
		return 5000 // Temperatura neutra-fría
	} else if ratio >= 0.5 {
		return 4000 // Temperatura neutra-cálida
	} else {
		return 3000 // Temperatura cálida
	}
}
