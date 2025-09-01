package main

import (
	"flag"
	"fyne.io/fyne/v2/app"
	"luznocturna/luz-nocturna/internal/controllers"
	"luznocturna/luz-nocturna/internal/views"
)

func main() {
	// Flags de línea de comandos
	trayMode := flag.Bool("tray", false, "Iniciar en modo bandeja del sistema")
	flag.Parse()

	// Crear la aplicación
	myApp := app.NewWithID("com.luznocturna.app")

	// Crear controlador
	controller := controllers.NewNightLightController()

	if *trayMode {
		// Modo bandeja del sistema (sin ventana visible)
		systrayManager := views.NewSystrayManager(controller, nil)
		systrayManager.Run() // Esto es bloqueante
	} else {
		// Modo ventana normal con soporte opcional de bandeja
		window := myApp.NewWindow("🌙 Luz Nocturna")
		window.CenterOnScreen()

		// Crear vista principal
		mainView := views.NewNightLightView(window, controller)

		// Crear manejador de bandeja (opcional, en segundo plano)
		go func() {
			systrayManager := views.NewSystrayManager(controller, mainView)
			systrayManager.Run()
		}()

		// Configurar comportamiento al cerrar
		window.SetCloseIntercept(func() {
			// En lugar de cerrar completamente, minimizar a bandeja
			window.Hide()
		})

		// Mostrar y ejecutar la aplicación
		window.ShowAndRun()
	}
}
