package styles

import (
	"fyne.io/fyne/v2/widget"
)

// Dimensiones de la aplicación
const (
	WindowWidth  = 320
	WindowHeight = 200

	// Padding y márgenes
	DefaultPadding = 20
	ElementSpacing = 10
	ButtonPadding  = 8

	// Tamaños de fuente
	TitleFontSize  = 16
	LabelFontSize  = 14
	ButtonFontSize = 12

	// Border radius (para futuras mejoras)
	BorderRadius = 12
	ButtonRadius = 8
)

// Función para aplicar estilos a botones
func StyleButton(btn *widget.Button, isPrimary bool) {
	if isPrimary {
		btn.Importance = widget.HighImportance
	} else {
		btn.Importance = widget.MediumImportance
	}
}

