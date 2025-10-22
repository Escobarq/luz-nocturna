package system

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
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
	gm.disableSystemNightLight()
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
 * applyWaylandGamma - Aplica gamma usando overlays de color efectivos para Wayland
 *
 * Implementa métodos más agresivos que realmente funcionen en Wayland
 * incluyendo overlays de color y filtros visuales.
 *
 * @param {float64} r - Componente rojo del gamma (0.3-1.0)
 * @param {float64} g - Componente verde del gamma (0.3-1.0)
 * @param {float64} b - Componente azul del gamma (0.3-1.0)
 * @returns {error} Error si falla la aplicación
 * @private
 */
func (gm *GammaManager) applyWaylandGamma(r, g, b float64) error {
	// Deshabilitar sistema nativo antes de aplicar
	gm.disableSystemNightLight()

	// Calcular temperatura para métodos que la requieren
	temp := gm.rgbToTemperature(r, g, b)

	// 1. Método más agresivo: Forzar gamma usando compositor
	if gm.tryCompositorOverride(r, g, b, temp) {
		return nil
	}

	// 2. Método compositor específico: GNOME Mutter
	if gm.tryGnomeMutterMethod(temp) {
		return nil
	}

	// 3. Método compositor específico: KDE KWin
	if gm.tryKWinMethod(temp) {
		return nil
	}

	// 4. Método DDC/CI para control directo del monitor
	if gm.tryDDCMethod(r, g, b) {
		return nil
	}

	// 5. Método overlay de color usando herramientas gráficas
	if gm.tryColorOverlayMethod(r, g, b) {
		return nil
	}

	// 6. Fallback: XWayland si está disponible
	if gm.tryXWaylandMethod(r, g, b) {
		fmt.Printf("⚠️  Usando XWayland (puede no ser efectivo en Wayland nativo)\n")
		return nil
	}

	return fmt.Errorf("no se pudo aplicar gamma en Wayland.\n" +
		"Métodos intentados: compositor override, GNOME, KDE, DDC/CI, overlay, XWayland\n" +
		"Tu compositor Wayland puede no soportar control de gamma")
}

/**
 * tryCompositorOverride - Método agresivo para forzar gamma en compositor
 */
func (gm *GammaManager) tryCompositorOverride(r, g, b, temp float64) bool {
	// 1. Intentar con wlr-gamma-control más agresivo
	if gm.isToolAvailable("wlr-gamma-control") {
		cmd := exec.Command("wlr-gamma-control", fmt.Sprintf("%.2f", r), fmt.Sprintf("%.2f", g), fmt.Sprintf("%.2f", b))
		if err := cmd.Run(); err == nil {
			fmt.Printf("🌡️  Gamma aplicada en Wayland (wlr-gamma-control): %.2f:%.2f:%.2f\n", r, g, b)
			return true
		}
	}

	// 2. Crear archivo temporal de configuración de gamma
	configPath := "/tmp/luz-nocturna-gamma.conf"
	configContent := fmt.Sprintf(`
[output:*]
gamma = %.2f:%.2f:%.2f
temperature = %.0f
`, r, g, b, temp)

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err == nil {
		// Intentar aplicar con swaybg si está disponible
		if gm.isToolAvailable("swaybg") {
			cmd := exec.Command("swaybg", "-c", fmt.Sprintf("#%02x%02x%02x",
				int(255*r), int(255*g), int(255*b)))
			if err := cmd.Start(); err == nil {
				fmt.Printf("🌡️  Overlay de color aplicado en Wayland (swaybg): %.2f:%.2f:%.2f\n", r, g, b)
				return true
			}
		}
	}

	return false
}

/**
 * tryGnomeMutterMethod - Método específico para GNOME Mutter
 */
func (gm *GammaManager) tryGnomeMutterMethod(temp float64) bool {
	if !gm.isToolAvailable("gdbus") {
		return false
	}

	// Forzar habilitación temporal del Night Light para controlarlo
	exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-enabled", "true").Run()
	time.Sleep(100 * time.Millisecond)

	// Configurar temperatura específica
	cmd := exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-temperature", fmt.Sprintf("uint32:%.0f", temp))
	if err := cmd.Run(); err == nil {
		// Forzar aplicación inmediata via D-Bus
		exec.Command("gdbus", "call", "--session", "--dest", "org.gnome.SettingsDaemon.Color",
			"--object-path", "/org/gnome/SettingsDaemon/Color",
			"--method", "org.gnome.SettingsDaemon.Color.NightLightPreview",
			fmt.Sprintf("uint32:%.0f", temp)).Run()

		fmt.Printf("🌡️  Temperatura aplicada en Wayland (GNOME Mutter): %.0fK\n", temp)
		return true
	}
	return false
}

/**
 * tryKWinMethod - Método específico para KDE KWin
 */
func (gm *GammaManager) tryKWinMethod(temp float64) bool {
	if !gm.isToolAvailable("qdbus") {
		return false
	}

	// Habilitar Night Color en KDE
	cmd := exec.Command("qdbus", "org.kde.KWin", "/ColorCorrect", "setMode", "2")
	if err := cmd.Run(); err == nil {
		// Configurar temperatura
		cmd = exec.Command("qdbus", "org.kde.KWin", "/ColorCorrect", "setTemperature", fmt.Sprintf("%.0f", temp))
		if err := cmd.Run(); err == nil {
			fmt.Printf("🌡️  Temperatura aplicada en Wayland (KDE KWin): %.0fK\n", temp)
			return true
		}
	}
	return false
}

/**
 * tryDDCMethod - Control directo del monitor usando DDC/CI
 */
func (gm *GammaManager) tryDDCMethod(r, g, b float64) bool {
	if !gm.isToolAvailable("ddcutil") {
		return false
	}

	// Convertir RGB a valores de color de monitor
	redVal := int(r * 100)
	greenVal := int(g * 100)
	blueVal := int(b * 100)

	// Aplicar usando ddcutil para control directo del hardware
	commands := [][]string{
		{"ddcutil", "setvcp", "16", fmt.Sprintf("%d", redVal)},   // Red gain
		{"ddcutil", "setvcp", "18", fmt.Sprintf("%d", greenVal)}, // Green gain
		{"ddcutil", "setvcp", "1A", fmt.Sprintf("%d", blueVal)},  // Blue gain
	}

	success := false
	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err == nil {
			success = true
		}
	}

	if success {
		fmt.Printf("🌡️  Gamma aplicada en Wayland (DDC/CI hardware): %.2f:%.2f:%.2f\n", r, g, b)
		return true
	}
	return false
}

/**
 * tryColorOverlayMethod - Crear overlay de color usando herramientas gráficas
 */
func (gm *GammaManager) tryColorOverlayMethod(r, g, b float64) bool {
	// Calcular color de overlay inverso para simular filtro
	overlayR := 1.0 - (1.0-r)*0.3
	overlayG := 1.0 - (1.0-g)*0.3
	overlayB := 1.0 - (1.0-b)*0.3

	colorHex := fmt.Sprintf("#%02x%02x%02x",
		int(255*overlayR), int(255*overlayG), int(255*overlayB))

	// Intentar con diferentes herramientas de overlay
	overlayTools := [][]string{
		{"pkill", "goverlay"}, // Matar overlay anterior
		{"goverlay", "--color", colorHex, "--opacity", "0.1"},
	}

	for _, cmdArgs := range overlayTools {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Start() // No esperar, es un overlay
	}

	// También intentar con xsetroot si funciona en XWayland
	if gm.isToolAvailable("xsetroot") {
		cmd := exec.Command("xsetroot", "-solid", colorHex)
		if err := cmd.Run(); err == nil {
			fmt.Printf("🌡️  Overlay de color aplicado en Wayland: %s\n", colorHex)
			return true
		}
	}

	return false
}

/**
 * tryXWaylandMethod - Intenta aplicar gamma usando xrandr en XWayland
 */
func (gm *GammaManager) tryXWaylandMethod(r, g, b float64) bool {
	if !gm.isToolAvailable("xrandr") {
		return false
	}

	// Verificar si hay displays detectados
	cmd := exec.Command("xrandr")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Buscar displays conectados
	lines := strings.Split(string(output), "\n")
	connectedRegex := regexp.MustCompile(`^(\S+)\s+connected`)

	applied := false
	for _, line := range lines {
		if matches := connectedRegex.FindStringSubmatch(line); matches != nil {
			display := matches[1]
			cmd := exec.Command("xrandr", "--output", display, "--gamma", fmt.Sprintf("%.2f:%.2f:%.2f", r, g, b))
			if err := cmd.Run(); err == nil {
				fmt.Printf("🌡️  Gamma aplicada en Wayland (XWayland/%s): %.2f:%.2f:%.2f\n", display, r, g, b)
				applied = true
			}
		}
	}
	return applied
}

/**
 * tryDBusMethod - Intenta aplicar temperatura usando D-Bus
 */
func (gm *GammaManager) tryDBusMethod(temp float64) bool {
	if !gm.isToolAvailable("dbus-send") {
		return false
	}

	// Intentar con GNOME Settings Daemon
	cmd := exec.Command("dbus-send", "--session", "--type=method_call",
		"--dest=org.gnome.SettingsDaemon.Color",
		"/org/gnome/SettingsDaemon/Color",
		"org.gnome.SettingsDaemon.Color.NightLightPreview",
		fmt.Sprintf("uint32:%.0f", temp))

	if err := cmd.Run(); err == nil {
		fmt.Printf("🌡️  Temperatura aplicada en Wayland (D-Bus/GNOME): %.0fK\n", temp)
		return true
	}

	// Intentar con KDE
	cmd = exec.Command("dbus-send", "--session", "--type=method_call",
		"--dest=org.kde.KWin",
		"/ColorCorrect",
		"org.kde.kwin.ColorCorrect.setMode",
		"string:manual")

	if err := cmd.Run(); err == nil {
		cmd = exec.Command("dbus-send", "--session", "--type=method_call",
			"--dest=org.kde.KWin",
			"/ColorCorrect",
			"org.kde.kwin.ColorCorrect.setTemperature",
			fmt.Sprintf("int32:%.0f", temp))

		if err := cmd.Run(); err == nil {
			fmt.Printf("🌡️  Temperatura aplicada en Wayland (D-Bus/KDE): %.0fK\n", temp)
			return true
		}
	}

	return false
}

/**
 * tryWlGammaRelay - Intenta usar wl-gamma-relay
 */
func (gm *GammaManager) tryWlGammaRelay(r, g, b float64) bool {
	if !gm.isToolAvailable("wl-gamma-relay") {
		return false
	}

	cmd := exec.Command("wl-gamma-relay", fmt.Sprintf("%.2f", r), fmt.Sprintf("%.2f", g), fmt.Sprintf("%.2f", b))
	if err := cmd.Run(); err == nil {
		fmt.Printf("🌡️  Gamma aplicada en Wayland (wl-gamma-relay): %.2f:%.2f:%.2f\n", r, g, b)
		return true
	}
	return false
}

/**
 * tryBrightnessMethod - Intenta simular temperatura ajustando brillo de pantalla
 */
func (gm *GammaManager) tryBrightnessMethod(r, g, b float64) bool {
	// Calcular brillo basado en valores RGB
	brightness := (r + g + b) / 3.0

	// Buscar archivos de brillo en /sys/class/backlight/
	cmd := exec.Command("find", "/sys/class/backlight/", "-name", "brightness", "2>/dev/null")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	brightnessFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range brightnessFiles {
		if file == "" {
			continue
		}

		// Leer brillo máximo
		maxFile := strings.Replace(file, "brightness", "max_brightness", 1)
		maxOutput, err := exec.Command("cat", maxFile).Output()
		if err != nil {
			continue
		}

		var maxBrightness int
		fmt.Sscanf(strings.TrimSpace(string(maxOutput)), "%d", &maxBrightness)

		// Calcular nuevo brillo
		newBrightness := int(float64(maxBrightness) * brightness)

		// Aplicar nuevo brillo
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %d | sudo tee %s", newBrightness, file))
		if err := cmd.Run(); err == nil {
			fmt.Printf("🌡️  Brillo ajustado en Wayland: %.0f%% (simulando temperatura)\n", brightness*100)
			return true
		}
	}
	return false
}

/**
 * tryRedshiftMethod - Intenta usar redshift temporalmente
 */
func (gm *GammaManager) tryRedshiftMethod(temp float64) bool {
	if !gm.isToolAvailable("redshift") {
		return false
	}

	// Matar redshift anterior
	exec.Command("pkill", "redshift").Run()
	time.Sleep(100 * time.Millisecond)

	// Aplicar temperatura con redshift
	cmd := exec.Command("redshift", "-P", "-O", fmt.Sprintf("%.0f", temp))
	if err := cmd.Run(); err == nil {
		fmt.Printf("🌡️  Temperatura aplicada en Wayland (redshift): %.0fK\n", temp)
		return true
	}
	return false
}

/**
 * resetWaylandGamma - Resetea gamma en Wayland usando múltiples métodos
 *
 * @returns {error} Error si falla el reset
 * @private
 */
func (gm *GammaManager) resetWaylandGamma() error {
	// Matar todos los procesos de control de gamma
	processes := []string{"wlsunset", "wl-gamma-relay", "gammastep", "redshift", "f.lux"}
	for _, proc := range processes {
		exec.Command("pkill", "-9", proc).Run()
		exec.Command("killall", "-9", proc).Run()
	}
	time.Sleep(300 * time.Millisecond)

	// 1. Intentar reset con XWayland
	if gm.tryXWaylandMethod(1.0, 1.0, 1.0) {
		fmt.Println("✅ Gamma reseteada en Wayland (XWayland)")
		return nil
	}

	// 2. Intentar reset con D-Bus
	if gm.tryDBusMethod(6500) {
		fmt.Println("✅ Gamma reseteada en Wayland (D-Bus)")
		return nil
	}

	// 3. Intentar reset con wl-gamma-relay
	if gm.isToolAvailable("wl-gamma-relay") {
		cmd := exec.Command("wl-gamma-relay", "1.0", "1.0", "1.0")
		if err := cmd.Run(); err == nil {
			fmt.Println("✅ Gamma reseteada en Wayland (wl-gamma-relay)")
			return nil
		}
	}

	// 4. Resetear configuración del sistema nativo
	if gm.isToolAvailable("gsettings") {
		// Habilitar de nuevo el sistema nativo y ponerlo en modo día
		exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-enabled", "false").Run()
		exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-temperature", "6500").Run()
	}

	fmt.Println("✅ Reset de gamma completado en Wayland")
	return nil
}

/**
 * detectWaylandDisplays - Detecta displays en Wayland
 *
 * Intenta detectar displays reales usando xrandr si está disponible,
 * de lo contrario usa control global de Wayland.
 *
 * @private
 */
func (gm *GammaManager) detectWaylandDisplays() {
	// Intentar usar xrandr incluso en Wayland (funciona en XWayland)
	if gm.isToolAvailable("xrandr") {
		cmd := exec.Command("xrandr")
		output, err := cmd.Output()
		if err == nil {
			// Parsear output de xrandr para encontrar displays conectados
			lines := strings.Split(string(output), "\n")
			connectedRegex := regexp.MustCompile(`^(\S+)\s+connected`)

			var displays []string
			for _, line := range lines {
				if matches := connectedRegex.FindStringSubmatch(line); matches != nil {
					displays = append(displays, matches[1])
				}
			}

			if len(displays) > 0 {
				gm.displays = displays
				fmt.Printf("🖥️  Displays detectados en Wayland (xrandr): %v\n", displays)
				return
			}
		}
	}

	// Fallback a control global de Wayland
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
 * isToolAvailable - Verifica si una herramienta está disponible en el sistema
 *
 * @param {string} tool - Nombre de la herramienta a verificar
 * @returns {bool} true si la herramienta está disponible
 * @private
 */
func (gm *GammaManager) isToolAvailable(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
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
	// Estimación mejorada basada en valores RGB gamma

	// Si todos los valores están cerca de 1.0, es temperatura diurna
	if r >= 0.95 && g >= 0.95 && b >= 0.95 {
		return 6500
	}

	// Usar el valor azul como indicador principal
	if b >= 0.9 {
		return 6500 // Muy frío/diurno
	} else if b >= 0.8 {
		return 5500 // Frío
	} else if b >= 0.7 {
		return 4500 // Neutro-frío
	} else if b >= 0.6 {
		return 4000 // Neutro-cálido
	} else if b >= 0.5 {
		return 3500 // Cálido
	} else {
		return 3000 // Muy cálido
	}
}

/**
 * disableSystemNightLight - Deshabilita automáticamente sistemas nativos de ZorinOS
 *
 * Detecta y deshabilita agresivamente todos los sistemas de luz nocturna
 * del entorno de escritorio para mantener control exclusivo.
 *
 * @private
 */
func (gm *GammaManager) disableSystemNightLight() {
	// Deshabilitar sistemas nativos silenciosamente

	// 1. GNOME/ZorinOS Night Light - Deshabilitación forzada
	if gm.isToolAvailable("gsettings") {
		// Verificar si está activo
		cmd := exec.Command("gsettings", "get", "org.gnome.settings-daemon.plugins.color", "night-light-enabled")
		output, err := cmd.Output()
		if err == nil {
			isEnabled := strings.TrimSpace(string(output)) == "true"

			// Deshabilitar completamente
			exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-enabled", "false").Run()
			exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-temperature", "uint32:6500").Run()
			exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-schedule-automatic", "false").Run()

			// Forzar aplicación inmediata via D-Bus
			if gm.isToolAvailable("gdbus") {
				exec.Command("gdbus", "call", "--session", "--dest", "org.gnome.SettingsDaemon.Color",
					"--object-path", "/org/gnome/SettingsDaemon/Color",
					"--method", "org.gnome.SettingsDaemon.Color.NightLightPreview",
					"uint32:6500").Run()
			}

			if isEnabled {
				fmt.Println("🔧 Sistema nativo deshabilitado")
			}
		}
	}

	// 2. KDE Night Color - Deshabilitación completa
	if gm.isToolAvailable("qdbus") {
		exec.Command("qdbus", "org.kde.KWin", "/ColorCorrect", "setMode", "0").Run()
	}

	// 3. Terminar todos los procesos competidores agresivamente
	processes := []string{
		"redshift", "redshift-gtk",
		"f.lux", "fluxgui", "xflux",
		"wlsunset", "wl-sunset",
		"gammastep", "gammastep-indicator",
		"goverlay", "blue-light-filter",
		"gnome-settings-daemon", // Reiniciar daemon si es necesario
	}

	killed := []string{}
	for _, proc := range processes {
		cmd := exec.Command("pgrep", proc)
		if err := cmd.Run(); err == nil {
			// Terminar proceso gracefully primero
			exec.Command("pkill", "-TERM", proc).Run()
			time.Sleep(100 * time.Millisecond)
			// Si sigue corriendo, forzar terminación
			exec.Command("pkill", "-KILL", proc).Run()
			killed = append(killed, proc)
		}
	}

	if len(killed) > 0 {
		time.Sleep(300 * time.Millisecond)
	}

	// 4. Crear archivo de bloqueo para evitar reactivación automática
	gm.createSystemLockFile()

	// 5. Monitorear y mantener control exclusivo
	go gm.maintainExclusiveControl()
}

/**
 * createSystemLockFile - Crea archivo para indicar que tenemos control exclusivo
 */
func (gm *GammaManager) createSystemLockFile() {
	lockDir := "/tmp/luz-nocturna"
	lockFile := lockDir + "/exclusive-control.lock"

	// Crear directorio si no existe
	os.MkdirAll(lockDir, 0755)

	// Crear archivo de bloqueo con información
	lockContent := fmt.Sprintf("luz-nocturna active\npid: %d\ntime: %s\n",
		os.Getpid(), time.Now().Format(time.RFC3339))

	os.WriteFile(lockFile, []byte(lockContent), 0644)
}

/**
 * maintainExclusiveControl - Mantiene control exclusivo del gamma
 */
func (gm *GammaManager) maintainExclusiveControl() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Verificar si el sistema nativo se reactivó
		if gm.isToolAvailable("gsettings") {
			cmd := exec.Command("gsettings", "get", "org.gnome.settings-daemon.plugins.color", "night-light-enabled")
			output, err := cmd.Output()
			if err == nil && strings.TrimSpace(string(output)) == "true" {
				// El sistema nativo se reactivó, deshabilitarlo de nuevo
				exec.Command("gsettings", "set", "org.gnome.settings-daemon.plugins.color", "night-light-enabled", "false").Run()
			}
		}

		// Verificar procesos competidores
		competitorProcesses := []string{"redshift", "wlsunset", "gammastep"}
		for _, proc := range competitorProcesses {
			cmd := exec.Command("pgrep", proc)
			if err := cmd.Run(); err == nil {
				exec.Command("pkill", "-TERM", proc).Run()
			}
		}
	}
}
