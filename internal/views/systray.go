package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"luznocturna/luz-nocturna/internal/controllers"
	"luznocturna/luz-nocturna/internal/models"
)

// SystrayManager - Manejador del icono de bandeja del sistema
type SystrayManager struct {
	controller *controllers.NightLightController
	mainView   *NightLightView
	app        fyne.App
}

// NewSystrayManager - Constructor del manejador de bandeja
func NewSystrayManager(app fyne.App, controller *controllers.NightLightController, mainView *NightLightView) *SystrayManager {
	return &SystrayManager{
		app:        app,
		controller: controller,
		mainView:   mainView,
	}
}

// CreateMenu - Crea y configura el menú de la bandeja del sistema
func (s *SystrayManager) CreateMenu() {
	if desk, ok := s.app.(desktop.App); ok {
		// 1. Crear el submenú de presets
		presetsSubMenu := fyne.NewMenu("Presets", // El título aquí es para la estructura interna
			fyne.NewMenuItem(fmt.Sprintf("🔥 Cálido (%.0fK)", models.CandleLightTemp), func() {
				s.applyTemperaturePreset(int(models.CandleLightTemp), "Cálido")
			}),
			fyne.NewMenuItem(fmt.Sprintf("🌅 Medio (%.0fK)", models.NeutralWhiteTemp), func() {
				s.applyTemperaturePreset(int(models.NeutralWhiteTemp), "Medio")
			}),
			fyne.NewMenuItem(fmt.Sprintf("☀️ Frío (%.0fK)", models.CoolWhiteTemp), func() {
				s.applyTemperaturePreset(int(models.CoolWhiteTemp), "Neutral")
			}),
			fyne.NewMenuItem(fmt.Sprintf("💡 Día (%.0fK)", models.DaylightTemp), func() {
				s.applyTemperaturePreset(int(models.DaylightTemp), "Día")
			}),
		)

		// 2. Crear el ítem de menú que contendrá el submenú
		presetsMenuItem := fyne.NewMenuItem("🌡️ Presets", nil)
		presetsMenuItem.ChildMenu = presetsSubMenu

		// 3. Crear el menú principal y añadir el ítem con el submenú
		menuItems := []*fyne.MenuItem{
			fyne.NewMenuItem("🌙 Aplicar", s.applyCurrentSettings),
			fyne.NewMenuItem("🔄 Resetear", s.resetToNormal),
			fyne.NewMenuItemSeparator(),
			presetsMenuItem, // Añadir el ítem que despliega el submenú
			fyne.NewMenuItemSeparator(),
		}

		if s.mainView != nil {
			menuItems = append(menuItems, fyne.NewMenuItem("📱 Mostrar", s.showMainWindow))
		}

		menuItems = append(menuItems, fyne.NewMenuItem("❌ Salir", func() {
			s.app.Quit()
		}))

		mainMenu := fyne.NewMenu("Luz Nocturna", menuItems...)

		desk.SetSystemTrayMenu(mainMenu)

		// Configurar icono
		iconData := GetOptimalIcon()
		if len(iconData) > 0 {
			desk.SetSystemTrayIcon(fyne.NewStaticResource("trayIcon", iconData))
		}
	}
}

func (s *SystrayManager) applyCurrentSettings() {
	_ = s.controller.ApplyNightLight()
}

func (s *SystrayManager) resetToNormal() {
	_ = s.controller.ResetNightLight()
}

func (s *SystrayManager) applyTemperaturePreset(temperature int, presetName string) {
	config := s.controller.GetConfig()
	config.Temperature = float64(temperature)

	_ = s.controller.ApplyNightLight()

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
