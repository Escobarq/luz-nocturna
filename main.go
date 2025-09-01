package main

import (
	"fyne.io/fyne/v2/app"

	"luznocturna/luz-nocturna/internal/controllers"
	"luznocturna/luz-nocturna/internal/views"
)

func main() {
	// Crear la aplicación
	myApp := app.NewWithID("com.luznocturna.app")

	// Crear ventana principal
	window := myApp.NewWindow("🌙 Luz Nocturna")
	window.CenterOnScreen()

	// Crear controlador
	controller := controllers.NewNightLightController()

	// Crear vista
	_ = views.NewNightLightView(window, controller)

	// Mostrar y ejecutar la aplicación
	window.ShowAndRun()
}
