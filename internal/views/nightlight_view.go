package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"luznocturna/luz-nocturna/internal/controllers"
	"luznocturna/luz-nocturna/internal/styles"
)

// NightLightView representa la vista principal de la aplicación
type NightLightView struct {
	controller        *controllers.NightLightController
	window            fyne.Window
	temperatureLabel  *widget.Label
	temperatureSlider *widget.Slider
	applyButton       *widget.Button
	resetButton       *widget.Button
}

// NewNightLightView crea una nueva vista
func NewNightLightView(window fyne.Window, controller *controllers.NightLightController) *NightLightView {
	view := &NightLightView{
		controller: controller,
		window:     window,
	}

	view.setupUI()
	return view
}

// setupUI configura todos los elementos de la interfaz
func (v *NightLightView) setupUI() {
	// Configurar ventana
	v.window.Resize(fyne.NewSize(styles.WindowWidth, styles.WindowHeight))
	v.window.SetFixedSize(true)

	// Crear widgets
	v.createWidgets()

	// Crear layout principal
	content := v.createMainLayout()

	// Establecer contenido
	v.window.SetContent(content)

	// Actualizar valores iniciales
	v.updateTemperatureDisplay()
}

// createWidgets crea todos los widgets de la interfaz
func (v *NightLightView) createWidgets() {
	config := v.controller.GetConfig()
	minTemp, maxTemp := v.controller.GetTemperatureRange()

	// Título con emoji
	title := widget.NewLabel("🌙 Luz Nocturna")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Label para mostrar temperatura actual
	v.temperatureLabel = widget.NewLabel("Temperatura de color: " + config.GetTemperatureString())
	v.temperatureLabel.Alignment = fyne.TextAlignCenter

	// Slider para temperatura
	v.temperatureSlider = widget.NewSlider(minTemp, maxTemp)
	v.temperatureSlider.Value = config.Temperature
	v.temperatureSlider.Step = 100
	v.temperatureSlider.OnChanged = v.onTemperatureChanged

	// Botón aplicar
	v.applyButton = widget.NewButton("Aplicar", v.onApplyClicked)
	styles.StyleButton(v.applyButton, true) // Botón primario

	// Botón reset
	v.resetButton = widget.NewButton("Reset", v.onResetClicked)
	styles.StyleButton(v.resetButton, false) // Botón secundario
}

// createMainLayout crea el layout principal de la aplicación
func (v *NightLightView) createMainLayout() fyne.CanvasObject {
	// Título
	title := widget.NewLabel("🌙 Luz Nocturna")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Contenedor del slider
	sliderContainer := container.NewVBox(
		v.temperatureLabel,
		v.temperatureSlider,
	)

	// Contenedor de botones
	buttonContainer := container.NewHBox(
		v.applyButton,
		v.resetButton,
	)

	// Layout principal
	mainContainer := container.NewVBox(
		title,
		widget.NewSeparator(),
		sliderContainer,
		buttonContainer,
	)

	// Contenedor con padding
	paddedContainer := container.NewPadded(mainContainer)

	return paddedContainer
}

// Eventos y callbacks

// onTemperatureChanged se ejecuta cuando cambia el valor del slider
func (v *NightLightView) onTemperatureChanged(value float64) {
	v.controller.UpdateTemperature(value)
	v.updateTemperatureDisplay()
}

// onApplyClicked se ejecuta cuando se presiona el botón Aplicar
func (v *NightLightView) onApplyClicked() {
	err := v.controller.ApplyNightLight()
	if err != nil {
		// Manejar error - por ahora solo lo mostramos en consola
		println("Error al aplicar luz nocturna:", err.Error())
		return
	}

	// Mostrar diálogo de confirmación
	v.showSuccessDialog("Luz nocturna aplicada correctamente")
}

// onResetClicked se ejecuta cuando se presiona el botón Reset
func (v *NightLightView) onResetClicked() {
	err := v.controller.ResetNightLight()
	if err != nil {
		println("Error al resetear luz nocturna:", err.Error())
		return
	}

	// Actualizar UI después del reset
	config := v.controller.GetConfig()
	v.temperatureSlider.Value = config.Temperature
	v.updateTemperatureDisplay()

	// Mostrar diálogo de confirmación
	v.showSuccessDialog("Configuración restablecida a valores por defecto")
}

// updateTemperatureDisplay actualiza el texto que muestra la temperatura actual
func (v *NightLightView) updateTemperatureDisplay() {
	config := v.controller.GetConfig()
	v.temperatureLabel.SetText("Temperatura de color: " + config.GetTemperatureString())
}

// showSuccessDialog muestra un diálogo de éxito
func (v *NightLightView) showSuccessDialog(message string) {
	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel(message),
			widget.NewButton("OK", func() {
				// Cerrar el diálogo - esto se implementará mejor más adelante
			}),
		),
		v.window.Canvas(),
	)
	dialog.Show()
}
