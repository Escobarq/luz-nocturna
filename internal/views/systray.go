package views

import (
	"fmt"
	"luznocturna/luz-nocturna/internal/controllers"

	"fyne.io/systray"
)

// SystrayManager - Manejador del icono de bandeja del sistema
type SystrayManager struct {
	controller      *controllers.NightLightController
	mainView        *NightLightView
	applyItem       *systray.MenuItem
	resetItem       *systray.MenuItem
	tempWarmItem    *systray.MenuItem
	tempNeutralItem *systray.MenuItem
	tempCoolItem    *systray.MenuItem
	tempDayItem     *systray.MenuItem
	showItem        *systray.MenuItem
	quitItem        *systray.MenuItem
}

// NewSystrayManager - Constructor del manejador de bandeja
func NewSystrayManager(controller *controllers.NightLightController, mainView *NightLightView) *SystrayManager {
	return &SystrayManager{
		controller: controller,
		mainView:   mainView,
	}
}

// Run - Ejecuta el bucle principal de la bandeja
func (s *SystrayManager) Run() {
	systray.Run(s.onReady, s.onExit)
}

// onReady - Callback ejecutado cuando la bandeja está lista
func (s *SystrayManager) onReady() {
	// Configurar icono del systray
	iconData := GetOptimalIcon()
	systray.SetIcon(iconData)
	systray.SetTitle("Luz Nocturna")
	systray.SetTooltip("Control de temperatura de color")

	s.applyItem = systray.AddMenuItem("🌙 Aplicar", "Aplica la temperatura actual")
	s.resetItem = systray.AddMenuItem("�� Resetear", "Restaura configuración normal")

	systray.AddSeparator()

	tempSubMenu := systray.AddMenuItem("🌡️ Presets", "Temperaturas predefinidas")
	s.tempWarmItem = tempSubMenu.AddSubMenuItem("🔥 Cálido (2700K)", "Temperatura cálida")
	s.tempNeutralItem = tempSubMenu.AddSubMenuItem("🌅 Medio (3500K)", "Temperatura media")
	s.tempCoolItem = tempSubMenu.AddSubMenuItem("☀️ Neutral (5000K)", "Temperatura neutra")
	s.tempDayItem = tempSubMenu.AddSubMenuItem("💡 Día (6500K)", "Temperatura día")

	systray.AddSeparator()

	if s.mainView != nil {
		s.showItem = systray.AddMenuItem("📱 Mostrar", "Mostrar ventana")
	}
	s.quitItem = systray.AddMenuItem("❌ Salir", "Salir")

	s.handleEvents()
}

// handleEvents - Maneja eventos del menú
func (s *SystrayManager) handleEvents() {
	go func() {
		for range s.applyItem.ClickedCh {
			s.applyCurrentSettings()
		}
	}()

	go func() {
		for range s.resetItem.ClickedCh {
			s.resetToNormal()
		}
	}()

	go func() {
		for range s.tempWarmItem.ClickedCh {
			s.applyTemperaturePreset(2700, "Cálido")
		}
	}()

	go func() {
		for range s.tempNeutralItem.ClickedCh {
			s.applyTemperaturePreset(3500, "Medio")
		}
	}()

	go func() {
		for range s.tempCoolItem.ClickedCh {
			s.applyTemperaturePreset(5000, "Neutral")
		}
	}()

	go func() {
		for range s.tempDayItem.ClickedCh {
			s.applyTemperaturePreset(6500, "Día")
		}
	}()

	if s.showItem != nil {
		go func() {
			for range s.showItem.ClickedCh {
				s.showMainWindow()
			}
		}()
	}

	go func() {
		for range s.quitItem.ClickedCh {
			systray.Quit()
		}
	}()
}

func (s *SystrayManager) applyCurrentSettings() {
	config := s.controller.GetConfig()
	err := s.controller.ApplyNightLight()
	if err != nil {
		systray.SetTooltip(fmt.Sprintf("Error: %v", err))
		return
	}
	systray.SetTooltip(fmt.Sprintf("Aplicado: %dK", int(config.Temperature)))
}

func (s *SystrayManager) resetToNormal() {
	err := s.controller.ResetNightLight()
	if err != nil {
		systray.SetTooltip(fmt.Sprintf("Error: %v", err))
		return
	}
	systray.SetTooltip("Reseteado a normal")
}

func (s *SystrayManager) applyTemperaturePreset(temperature int, presetName string) {
	config := s.controller.GetConfig()
	config.Temperature = float64(temperature)

	err := s.controller.ApplyNightLight()
	if err != nil {
		systray.SetTooltip(fmt.Sprintf("Error: %v", err))
		return
	}

	systray.SetTooltip(fmt.Sprintf("%s (%dK) aplicado", presetName, temperature))

	if s.mainView != nil {
		s.mainView.updateTemperatureDisplay()
	}
}

func (s *SystrayManager) showMainWindow() {
	if s.mainView != nil && s.mainView.window != nil {
		s.mainView.window.Show()
		s.mainView.window.RequestFocus()
	}
}

func (s *SystrayManager) onExit() {
	// Limpieza si es necesaria
}
