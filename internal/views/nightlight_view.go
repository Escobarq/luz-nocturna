package views

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"luznocturna/luz-nocturna/internal/controllers"
	"luznocturna/luz-nocturna/internal/models"
	"luznocturna/luz-nocturna/internal/styles"
)

/**
 * NightLightView - Vista principal de la aplicación de luz nocturna
 *
 * Maneja toda la interfaz gráfica principal incluyendo controles de temperatura,
 * presets, botones de acción e información del sistema. Implementa el patrón MVC
 * como la capa de Vista que interactúa con el usuario.
 *
 * @struct {NightLightView}
 * @property {*controllers.NightLightController} controller - Controlador principal
 * @property {fyne.Window} window - Ventana principal de la aplicación
 * @property {*widget.Label} temperatureLabel - Etiqueta que muestra temperatura actual
 * @property {*widget.Slider} temperatureSlider - Control deslizante de temperatura
 * @property {*widget.Label} presetLabel - Etiqueta que muestra el preset actual
 * @property {*widget.Button} applyButton - Botón para aplicar configuración
 * @property {*widget.Button} resetButton - Botón para resetear a valores normales
 * @property {*widget.Button} toggleButton - Botón para alternar on/off
 * @property {*widget.Label} displayInfo - Información de displays detectados
 * @property {*fyne.Container} presetButtons - Contenedor de botones de presets
 */
type NightLightView struct {
	controller        *controllers.NightLightController
	window            fyne.Window
	temperatureLabel  *widget.Label
	temperatureSlider *widget.Slider
	presetLabel       *widget.Label
	applyButton       *widget.Button
	resetButton       *widget.Button
	toggleButton      *widget.Button
	displayInfo       *widget.Label
	presetButtons     *fyne.Container
	scheduleCheck     *widget.Check
	startTimeEntry    *widget.Entry
	endTimeEntry      *widget.Entry
	nightTempSlider   *widget.Slider
	dayTempSlider     *widget.Slider
	transitionSlider  *widget.Slider
	scheduleInfo      *widget.Label
}

/**
 * NewNightLightView - Constructor de la vista principal
 *
 * Crea una nueva instancia de la vista principal de la aplicación,
 * inicializa todos los componentes de la interfaz y configura los eventos.
 *
 * @param {fyne.Window} window - Ventana donde se mostrará la vista
 * @param {*controllers.NightLightController} controller - Controlador principal
 * @returns {*NightLightView} Nueva instancia de la vista
 * @example
 *   window := app.NewWindow("Luz Nocturna")
 *   controller := controllers.NewNightLightController()
 *   view := NewNightLightView(window, controller)
 */
func NewNightLightView(window fyne.Window, controller *controllers.NightLightController) *NightLightView {
	view := &NightLightView{
		controller: controller,
		window:     window,
	}

	view.setupUI()
	return view
}

/**
 * setupUI - Configura todos los elementos de la interfaz
 *
 * Método privado que inicializa la interfaz gráfica completa:
 * - Configura el tamaño y propiedades de la ventana
 * - Crea todos los widgets necesarios
 * - Establece el layout principal
 * - Actualiza valores iniciales
 *
 * @private
 */
func (v *NightLightView) setupUI() {
	// Configurar ventana principal
	v.window.Resize(fyne.NewSize(styles.WindowWidth, styles.WindowHeight+200))
	v.window.SetFixedSize(false)

	// Crear todos los widgets de la interfaz
	v.createWidgets()

	// Crear y establecer el layout principal
	content := v.createMainLayout()
	v.window.SetContent(content)

	// Sincronizar estado inicial con el modelo
	v.updateTemperatureDisplay()
	v.updateDisplayInfo()

	// Iniciar actualizador de información de programación
	v.startScheduleInfoUpdater()
}

/**
 * createWidgets - Crea todos los widgets de la interfaz
 *
 * Inicializa todos los componentes de la UI incluyendo labels, sliders,
 * botones y contenedores. Configura eventos y valores iniciales.
 *
 * @private
 */
func (v *NightLightView) createWidgets() {
	config := v.controller.GetConfig()
	minTemp, maxTemp := v.controller.GetTemperatureRange()

	// === LABELS DE INFORMACIÓN ===
	v.temperatureLabel = widget.NewLabel("Temperatura de color: " + config.GetTemperatureString())
	v.temperatureLabel.Alignment = fyne.TextAlignCenter

	v.presetLabel = widget.NewLabel(models.Presets.GetPresetName(config.Temperature))
	v.presetLabel.Alignment = fyne.TextAlignCenter
	v.presetLabel.TextStyle = fyne.TextStyle{Italic: true}

	// === CONTROL DESLIZANTE ===
	v.temperatureSlider = widget.NewSlider(minTemp, maxTemp)
	v.temperatureSlider.Value = config.Temperature
	v.temperatureSlider.Step = 100
	v.temperatureSlider.OnChanged = v.onTemperatureChanged

	// === BOTONES DE PRESETS ===
	v.createPresetButtons()

	// === BOTONES PRINCIPALES ===
	v.applyButton = widget.NewButton("🔥 Aplicar", v.onApplyClicked)
	styles.StyleButton(v.applyButton, true) // Botón primario

	v.resetButton = widget.NewButton("↺ Reset", v.onResetClicked)
	styles.StyleButton(v.resetButton, false) // Botón secundario

	v.toggleButton = widget.NewButton("🔄 Toggle", v.onToggleClicked)
	styles.StyleButton(v.toggleButton, false)

	// === INFORMACIÓN DEL SISTEMA ===
	displays := v.controller.GetDisplays()
	v.displayInfo = widget.NewLabel(fmt.Sprintf("📺 Displays: %v", displays))
	v.displayInfo.TextStyle = fyne.TextStyle{Monospace: true}

	// === CONTROLES DE PROGRAMACIÓN AUTOMÁTICA ===
	v.createScheduleWidgets()
}

/**
 * createScheduleWidgets - Crea los controles de programación automática
 *
 * @private
 */
func (v *NightLightView) createScheduleWidgets() {
	schedule := v.controller.GetScheduleConfig()

	// Checkbox para habilitar/deshabilitar programación
	v.scheduleCheck = widget.NewCheck("🕐 Programación automática", v.onScheduleToggled)
	v.scheduleCheck.SetChecked(v.controller.IsScheduleEnabled())

	// Entradas de tiempo
	v.startTimeEntry = widget.NewEntry()
	v.startTimeEntry.SetText(schedule.StartTime)
	v.startTimeEntry.OnChanged = v.onScheduleTimeChanged

	v.endTimeEntry = widget.NewEntry()
	v.endTimeEntry.SetText(schedule.EndTime)
	v.endTimeEntry.OnChanged = v.onScheduleTimeChanged

	// Sliders de temperatura
	v.nightTempSlider = widget.NewSlider(3000, 6500)
	v.nightTempSlider.Value = schedule.NightTemp
	v.nightTempSlider.Step = 100
	v.nightTempSlider.OnChanged = v.onScheduleTempChanged

	v.dayTempSlider = widget.NewSlider(3000, 6500)
	v.dayTempSlider.Value = schedule.DayTemp
	v.dayTempSlider.Step = 100
	v.dayTempSlider.OnChanged = v.onScheduleTempChanged

	// Slider de tiempo de transición
	v.transitionSlider = widget.NewSlider(0, 60)
	v.transitionSlider.Value = float64(schedule.TransitionTime)
	v.transitionSlider.Step = 5
	v.transitionSlider.OnChanged = v.onScheduleTempChanged

	// Información de próximo cambio
	v.scheduleInfo = widget.NewLabel("Programación deshabilitada")
	v.scheduleInfo.TextStyle = fyne.TextStyle{Italic: true}

	v.updateScheduleInfo()
}

/**
 * createPresetButtons - Crea los botones de presets de temperatura
 *
 * Genera botones rápidos para temperaturas predefinidas comunes:
 * Cálida (3000K), Neutra (4500K), Fría (5500K), Diurna (6500K)
 *
 * @private
 */
func (v *NightLightView) createPresetButtons() {
	presets := []struct {
		name string
		temp float64
		icon string
	}{
		{"Cálida", models.CandleLightTemp, "🕯️"},
		{"Neutra", models.NeutralWhiteTemp, "☀️"},
		{"Fría", models.CoolWhiteTemp, "🌤️"},
		{"Diurna", models.DaylightTemp, "☀️"},
	}

	var buttons []fyne.CanvasObject
	for _, preset := range presets {
		temp := preset.temp // Capturar valor para closure
		btn := widget.NewButton(preset.icon+" "+preset.name, func() {
			v.controller.UpdateTemperature(temp)
			v.temperatureSlider.Value = temp
			v.updateTemperatureDisplay()
		})
		buttons = append(buttons, btn)
	}

	v.presetButtons = container.NewGridWithColumns(2, buttons...)
}

/**
 * createMainLayout - Crea el layout principal de la aplicación
 *
 * Organiza todos los widgets en un layout vertical con separadores,
 * creando una interfaz limpia y bien organizada.
 *
 * @returns {fyne.CanvasObject} Contenedor principal listo para mostrar
 * @private
 */
func (v *NightLightView) createMainLayout() fyne.CanvasObject {
	// Título principal con emoji
	title := widget.NewLabel("🌙 Luz Nocturna")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Sección de control de temperatura
	tempContainer := container.NewVBox(
		v.temperatureLabel,
		v.presetLabel,
		v.temperatureSlider,
	)

	// Sección de presets rápidos
	presetSection := container.NewVBox(
		widget.NewLabel("🎨 Presets Rápidos:"),
		v.presetButtons,
	)

	// Botones principales de acción
	buttonContainer := container.NewGridWithColumns(3,
		v.applyButton,
		v.resetButton,
		v.toggleButton,
	)

	// Sección de programación automática
	scheduleSection := v.createScheduleSection()

	// Layout principal con separadores para claridad visual
	mainContainer := container.NewVBox(
		title,
		widget.NewSeparator(),
		tempContainer,
		widget.NewSeparator(),
		presetSection,
		widget.NewSeparator(),
		buttonContainer,
		widget.NewSeparator(),
		scheduleSection,
		widget.NewSeparator(),
		v.displayInfo,
	)

	// Contenedor con padding para mejor apariencia
	return container.NewPadded(mainContainer)
}

/**
 * createScheduleSection - Crea la sección de programación automática
 *
 * @returns {fyne.CanvasObject} Contenedor de la sección de programación
 * @private
 */
func (v *NightLightView) createScheduleSection() fyne.CanvasObject {
	// Contenedor principal de programación
	scheduleContainer := container.NewVBox(
		v.scheduleCheck,
	)

	// Controles de horarios (solo se muestran si está habilitado)
	timeContainer := container.NewGridWithColumns(4,
		widget.NewLabel("Inicio:"),
		v.startTimeEntry,
		widget.NewLabel("Fin:"),
		v.endTimeEntry,
	)

	// Controles de temperatura
	tempContainer := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("🌙 Temperatura nocturna: %.0fK", v.nightTempSlider.Value)),
		v.nightTempSlider,
		widget.NewLabel(fmt.Sprintf("☀️ Temperatura diurna: %.0fK", v.dayTempSlider.Value)),
		v.dayTempSlider,
	)

	// Control de transición
	transitionContainer := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("⏱️ Transición: %.0f min", v.transitionSlider.Value)),
		v.transitionSlider,
	)

	// Información de estado
	infoContainer := container.NewVBox(
		v.scheduleInfo,
	)

	// Crear contenedor colapsable para controles de programación
	configContainer := container.NewVBox()

	// Agregar controles condicionalmente
	if v.controller.IsScheduleEnabled() {
		configContainer.Add(timeContainer)
		configContainer.Add(tempContainer)
		configContainer.Add(transitionContainer)
	}

	scheduleContainer.Add(configContainer)
	scheduleContainer.Add(infoContainer)

	return container.NewVBox(
		widget.NewLabel("🕐 Programación Automática:"),
		scheduleContainer,
	)
}

// =====================================================
// MANEJADORES DE EVENTOS (Event Handlers)
// =====================================================

/**
 * onTemperatureChanged - Manejador de cambio en el slider de temperatura
 *
 * Se ejecuta cuando el usuario mueve el slider. Actualiza el modelo
 * y la interfaz en tiempo real para mostrar el cambio.
 *
 * @param {float64} value - Nueva temperatura seleccionada en Kelvin
 * @callback - Evento del slider
 */
func (v *NightLightView) onTemperatureChanged(value float64) {
	v.controller.UpdateTemperature(value)
	v.updateTemperatureDisplay()
}

/**
 * onApplyClicked - Manejador del botón Aplicar
 *
 * Aplica la temperatura actual al sistema usando el controlador.
 * Muestra feedback visual del resultado (éxito o error).
 *
 * @callback - Evento del botón Aplicar
 */
func (v *NightLightView) onApplyClicked() {
	err := v.controller.ApplyNightLight()
	if err != nil {
		v.showErrorDialog("❌ Error al aplicar", err.Error())
		return
	}

	config := v.controller.GetConfig()
	message := fmt.Sprintf("🌡️ Aplicada: %s", config.GetTemperatureString())
	v.showSuccessDialog(message)
}

/**
 * onResetClicked - Manejador del botón Reset
 *
 * Resetea la configuración a valores normales (6500K) y actualiza
 * tanto el sistema como la interfaz.
 *
 * @callback - Evento del botón Reset
 */
func (v *NightLightView) onResetClicked() {
	err := v.controller.ResetNightLight()
	if err != nil {
		v.showErrorDialog("❌ Error al resetear", err.Error())
		return
	}

	// Actualizar UI después del reset
	config := v.controller.GetConfig()
	v.temperatureSlider.Value = config.Temperature
	v.updateTemperatureDisplay()

	v.showSuccessDialog("✅ Gamma reseteada a valores normales")
}

/**
 * onScheduleToggled - Manejador del checkbox de programación automática
 *
 * @param {bool} enabled - Estado del checkbox
 * @callback - Evento del checkbox
 */
func (v *NightLightView) onScheduleToggled(enabled bool) {
	v.controller.EnableSchedule(enabled)
	v.refreshScheduleSection()
	v.updateScheduleInfo()
}

/**
 * onScheduleTimeChanged - Manejador de cambios en entradas de tiempo
 *
 * @param {string} text - Nuevo texto en la entrada
 * @callback - Evento de cambio en entradas de tiempo
 */
func (v *NightLightView) onScheduleTimeChanged(text string) {
	if !v.controller.IsScheduleEnabled() {
		return
	}

	v.updateScheduleConfiguration()
}

/**
 * onScheduleTempChanged - Manejador de cambios en sliders de temperatura
 *
 * @param {float64} value - Nuevo valor del slider
 * @callback - Evento de cambio en sliders
 */
func (v *NightLightView) onScheduleTempChanged(value float64) {
	if !v.controller.IsScheduleEnabled() {
		return
	}

	v.updateScheduleConfiguration()
	v.refreshScheduleSection() // Actualizar labels de temperatura
}

/**
 * updateScheduleConfiguration - Actualiza la configuración de horarios
 *
 * @private
 */
func (v *NightLightView) updateScheduleConfiguration() {
	// Obtener valores actuales de la UI
	startTime := v.startTimeEntry.Text
	endTime := v.endTimeEntry.Text
	nightTemp := v.nightTempSlider.Value
	dayTemp := v.dayTempSlider.Value
	transitionTime := int(v.transitionSlider.Value)

	// Actualizar configuración
	v.controller.UpdateScheduleConfig(startTime, endTime, nightTemp, dayTemp, transitionTime)

	// Actualizar información
	v.updateScheduleInfo()
}

/**
 * onToggleClicked - Manejador del botón Toggle
 *
 * Alterna entre activar y desactivar la luz nocturna.
 * Si está activa la desactiva, si está inactiva la activa.
 *
 * @callback - Evento del botón Toggle
 */
func (v *NightLightView) onToggleClicked() {
	err := v.controller.ToggleNightLight()
	if err != nil {
		v.showErrorDialog("❌ Error al cambiar estado", err.Error())
		return
	}

	config := v.controller.GetConfig()
	var message string
	if config.IsActive {
		message = "🔥 Luz nocturna activada"
	} else {
		message = "❄️ Luz nocturna desactivada"
	}

	// Actualizar UI
	v.temperatureSlider.Value = config.Temperature
	v.updateTemperatureDisplay()
	v.showSuccessDialog(message)
}

// =====================================================
// MÉTODOS DE ACTUALIZACIÓN DE UI
// =====================================================

/**
 * updateTemperatureDisplay - Actualiza la visualización de temperatura
 *
 * Sincroniza los labels de temperatura y preset con el estado actual
 * del modelo. Se llama cada vez que cambia la temperatura.
 *
 * @private
 */
func (v *NightLightView) updateTemperatureDisplay() {
	config := v.controller.GetConfig()
	v.temperatureLabel.SetText("🌡️ Temperatura: " + config.GetTemperatureString())
	v.presetLabel.SetText("✨ " + models.Presets.GetPresetName(config.Temperature))
}

/**
 * updateDisplayInfo - Actualiza la información de displays
 *
 * Refresca la información de displays detectados por el sistema.
 * Útil cuando se conectan/desconectan monitores.
 *
 * @private
 */
func (v *NightLightView) updateDisplayInfo() {
	displays := v.controller.GetDisplays()
	v.displayInfo.SetText(fmt.Sprintf("📺 Displays: %v", displays))
}

/**
 * updateScheduleInfo - Actualiza la información de programación automática
 *
 * @private
 */
func (v *NightLightView) updateScheduleInfo() {
	if !v.controller.IsScheduleEnabled() {
		v.scheduleInfo.SetText("Programación deshabilitada")
		return
	}

	description, temp, duration := v.controller.GetNextScheduleChange()

	if duration > 0 {
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		v.scheduleInfo.SetText(fmt.Sprintf("🔔 %s en %02d:%02d (%.0fK)",
			description, hours, minutes, temp))
	} else {
		v.scheduleInfo.SetText("🔔 " + description)
	}
}

/**
 * updateScheduleLabels - Actualiza los labels de los sliders de programación
 *
 * @private
 */
func (v *NightLightView) updateScheduleLabels() {
	// Esta función se llamará desde createScheduleSection cuando se recree el layout
	// Los labels se actualizan automáticamente en createScheduleSection
}

/**
 * refreshScheduleSection - Refresca la sección de programación automática
 *
 * @private
 */
func (v *NightLightView) refreshScheduleSection() {
	// Ajustar tamaño de ventana según estado de programación
	if v.controller.IsScheduleEnabled() {
		v.window.Resize(fyne.NewSize(styles.WindowWidth, styles.WindowHeight+300))
	} else {
		v.window.Resize(fyne.NewSize(styles.WindowWidth, styles.WindowHeight+150))
	}

	// Recrear el contenido de la ventana para mostrar/ocultar controles de programación
	content := v.createMainLayout()
	v.window.SetContent(content)
}

/**
 * startScheduleInfoUpdater - Inicia el actualizador automático de información de programación
 *
 * @private
 */
func (v *NightLightView) startScheduleInfoUpdater() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if v.controller.IsScheduleEnabled() {
				v.updateScheduleInfo()
			}
		}
	}()
}

// =====================================================
// SISTEMA DE DIÁLOGOS
// =====================================================

/**
 * showSuccessDialog - Muestra un diálogo de éxito auto-cerrable
 *
 * Presenta un mensaje de confirmación que se cierra automáticamente
 * después de 2 segundos para no interrumpir el flujo de trabajo.
 *
 * @param {string} message - Mensaje a mostrar al usuario
 * @example
 *   v.showSuccessDialog("✅ Configuración aplicada")
 */
func (v *NightLightView) showSuccessDialog(message string) {
	info := dialog.NewInformation("✅ Éxito", message, v.window)
	info.Show()

	// Auto-cerrar después de 2 segundos
	go func() {
		time.Sleep(2 * time.Second)
		info.Hide()
	}()
}

/**
 * showErrorDialog - Muestra un diálogo de error
 *
 * Presenta un error al usuario de forma clara. No se auto-cierra
 * para que el usuario pueda leer el mensaje completo.
 *
 * @param {string} title - Título del diálogo de error
 * @param {string} message - Mensaje detallado del error
 * @example
 *   v.showErrorDialog("Error de sistema", "No se pudo conectar al display")
 */
func (v *NightLightView) showErrorDialog(title, message string) {
	dialog.ShowError(fmt.Errorf("%s: %s", title, message), v.window)
}
