package views

import (
	"fyne.io/systray"
	"luznocturna/luz-nocturna/internal/controllers"
	"luznocturna/luz-nocturna/internal/models"
)

/**
 * SystrayManager - Manejador del icono de bandeja del sistema
 * 
 * Controla el icono de la bandeja del sistema (system tray) proporcionando
 * acceso rápido a las funcionalidades principales sin abrir la ventana principal.
 * Incluye menú contextual con presets, controles y información del estado.
 * 
 * @struct {SystrayManager}
 * @property {*controllers.NightLightController} controller - Controlador principal
 * @property {*NightLightView} mainView - Referencia a la vista principal (puede ser nil)
 * @property {*systray.MenuItem} applyItem - Item del menú para aplicar
 * @property {*systray.MenuItem} resetItem - Item del menú para resetear
 * @property {*systray.MenuItem} tempWarmItem - Preset de temperatura cálida
 * @property {*systray.MenuItem} tempNeutralItem - Preset de temperatura neutra
 * @property {*systray.MenuItem} tempCoolItem - Preset de temperatura fría
 * @property {*systray.MenuItem} tempDayItem - Preset de temperatura diurna
 * @property {*systray.MenuItem} showItem - Item para mostrar ventana principal
 * @property {*systray.MenuItem} quitItem - Item para salir de la aplicación
 */
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

/**
 * NewSystrayManager - Constructor del manejador de bandeja
 * 
 * Crea una nueva instancia del manejador de bandeja del sistema.
 * La mainView puede ser nil si se ejecuta en modo solo-bandeja.
 * 
 * @param {*controllers.NightLightController} controller - Controlador principal
 * @param {*NightLightView} mainView - Vista principal (opcional, puede ser nil)
 * @returns {*SystrayManager} Nueva instancia del manejador
 * @example
 *   manager := NewSystrayManager(controller, mainView)
 *   manager.Start() // Iniciar bandeja (método bloqueante)
 */
func NewSystrayManager(controller *controllers.NightLightController, mainView *NightLightView) *SystrayManager {
	return &SystrayManager{
		controller: controller,
		mainView:   mainView,
	}
}

/**
 * Start - Inicia el sistema de bandeja del sistema
 * 
 * Método bloqueante que inicializa y mantiene activa la bandeja del sistema.
 * Debe ejecutarse en el hilo principal o como goroutine según el uso.
 * 
 * @example
 *   // Modo bloqueante (solo bandeja)
 *   manager.Start()
 *   
 *   // Como goroutine (con ventana principal)
 *   go manager.Start()
 */
func (s *SystrayManager) Start() {
	systray.Run(s.onReady, s.onExit)
}

/**
 * onReady - Callback ejecutado cuando la bandeja está lista
 * 
 * Inicializa el icono, tooltip, elementos del menú y configura
 * los manejadores de eventos. Se ejecuta automáticamente.
 * 
 * @callback - Evento interno de systray
 * @private
 */
func (s *SystrayManager) onReady() {
	// Configurar icono y información básica
	systray.SetIcon(GetOptimalIcon()) // Usar icono optimizado según el sistema
	systray.SetTitle("🌙 Luz Nocturna")
	systray.SetTooltip("Control de temperatura de color - Clic para opciones")
	
	// Crear estructura del menú
	s.createMenuItems()
	
	// Configurar eventos de todos los elementos
	s.setupEventHandlers()
	
	// Sincronizar estado inicial
	s.updateMenuState()
}

/**
 * onExit - Callback ejecutado al cerrar la bandeja
 * 
 * Limpia recursos y realiza tareas de finalización.
 * Se ejecuta automáticamente al salir.
 * 
 * @callback - Evento interno de systray
 * @private
 */
func (s *SystrayManager) onExit() {
	// Aquí podrías agregar limpieza de recursos si fuera necesario
	// Por ejemplo: cerrar conexiones, guardar configuración final, etc.
}

/**
 * createMenuItems - Crea todos los elementos del menú contextual
 * 
 * Construye la estructura completa del menú de la bandeja incluyendo:
 * - Información de estado actual
 * - Presets de temperatura rápidos
 * - Acciones principales (aplicar, reset)
 * - Controles de ventana y aplicación
 * 
 * @private
 */
func (s *SystrayManager) createMenuItems() {
	config := s.controller.GetConfig()
	
	// === INFORMACIÓN DE ESTADO ===
	statusItem := systray.AddMenuItem("🌡️ "+config.GetTemperatureString(), "Temperatura actual")
	statusItem.Disable() // Solo informativo
	
	systray.AddSeparator()
	
	// === PRESETS DE TEMPERATURA ===
	tempMenu := systray.AddMenuItem("🎨 Presets de Temperatura", "Seleccionar temperatura predefinida")
	s.tempWarmItem = tempMenu.AddSubMenuItem("🕯️ Cálida (3000K)", "Temperatura cálida para la noche")
	s.tempNeutralItem = tempMenu.AddSubMenuItem("☀️ Neutra (4500K)", "Temperatura equilibrada")
	s.tempCoolItem = tempMenu.AddSubMenuItem("🌤️ Fría (5500K)", "Temperatura fría para el día")
	s.tempDayItem = tempMenu.AddSubMenuItem("☀️ Diurna (6500K)", "Temperatura de luz natural")
	
	systray.AddSeparator()
	
	// === ACCIONES PRINCIPALES ===
	s.applyItem = systray.AddMenuItem("🔥 Aplicar", "Aplicar temperatura actual al sistema")
	s.resetItem = systray.AddMenuItem("↺ Reset", "Resetear a valores normales (6500K)")
	
	systray.AddSeparator()
	
	// === CONTROLES DE APLICACIÓN ===
	if s.mainView != nil {
		s.showItem = systray.AddMenuItem("📱 Mostrar Ventana", "Abrir ventana principal")
		systray.AddSeparator()
	}
	
	// Información del sistema
	infoItem := systray.AddMenuItem("ℹ️ Información", "Información del sistema")
	infoItem.Disable()
	
	systray.AddSeparator()
	
	// === SALIR ===
	s.quitItem = systray.AddMenuItem("❌ Salir", "Cerrar aplicación completamente")
}

/**
 * setupEventHandlers - Configura todos los manejadores de eventos del menú
 * 
 * Establece goroutines que escuchan los clicks en cada elemento del menú
 * y ejecutan las acciones correspondientes. Maneja concurrencia de forma segura.
 * 
 * @private
 */
func (s *SystrayManager) setupEventHandlers() {
	// === MANEJADORES DE PRESETS ===
	go func() {
		for {
			select {
			case <-s.tempWarmItem.ClickedCh:
				s.setTemperaturePreset(models.CandleLightTemp, "🕯️ Cálida")
			case <-s.tempNeutralItem.ClickedCh:
				s.setTemperaturePreset(models.NeutralWhiteTemp, "☀️ Neutra")
			case <-s.tempCoolItem.ClickedCh:
				s.setTemperaturePreset(models.CoolWhiteTemp, "🌤️ Fría")
			case <-s.tempDayItem.ClickedCh:
				s.setTemperaturePreset(models.DaylightTemp, "☀️ Diurna")
			}
		}
	}()
	
	// === MANEJADORES DE ACCIONES PRINCIPALES ===
	go func() {
		for {
			select {
			case <-s.applyItem.ClickedCh:
				s.handleApplyAction()
			case <-s.resetItem.ClickedCh:
				s.handleResetAction()
			}
		}
	}()
	
	// === MANEJADORES DE CONTROLES DE APLICACIÓN ===
	go func() {
		for {
			select {
			case <-s.showItem.ClickedCh:
				if s.mainView != nil {
					s.showMainWindow()
				}
			case <-s.quitItem.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

/**
 * setTemperaturePreset - Aplica un preset de temperatura específico
 * 
 * Actualiza la temperatura en el controlador y sincroniza con la vista principal
 * si está disponible. También actualiza el estado del menú de bandeja.
 * 
 * @param {float64} temp - Temperatura en Kelvin del preset
 * @param {string} name - Nombre descriptivo del preset para logging
 * @private
 */
func (s *SystrayManager) setTemperaturePreset(temp float64, name string) {
	s.controller.UpdateTemperature(temp)
	
	// Sincronizar con vista principal si existe
	if s.mainView != nil {
		s.mainView.temperatureSlider.Value = temp
		s.mainView.updateTemperatureDisplay()
	}
	
	// Actualizar tooltip de la bandeja
	s.updateMenuState()
	
	// Log para debug (podrías usar un logger más sofisticado)
	println("📊 Preset aplicado desde bandeja:", name, "->", fmt.Sprintf("%.0fK", temp))
}

/**
 * handleApplyAction - Maneja la acción de aplicar temperatura
 * 
 * Ejecuta la aplicación de temperatura al sistema y actualiza el estado.
 * Maneja errores de forma silenciosa ya que no hay UI visible.
 * 
 * @private
 */
func (s *SystrayManager) handleApplyAction() {
	err := s.controller.ApplyNightLight()
	if err != nil {
		println("⚠️ Error al aplicar desde bandeja:", err.Error())
		return
	}
	
	s.updateMenuState()
	println("✅ Temperatura aplicada desde bandeja")
}

/**
 * handleResetAction - Maneja la acción de reset
 * 
 * Resetea la configuración a valores normales y actualiza la UI.
 * 
 * @private
 */
func (s *SystrayManager) handleResetAction() {
	err := s.controller.ResetNightLight()
	if err != nil {
		println("⚠️ Error al resetear desde bandeja:", err.Error())
		return
	}
	
	// Sincronizar con vista principal si existe
	if s.mainView != nil {
		config := s.controller.GetConfig()
		s.mainView.temperatureSlider.Value = config.Temperature
		s.mainView.updateTemperatureDisplay()
	}
	
	s.updateMenuState()
	println("↺ Reseteado desde bandeja")
}

/**
 * updateMenuState - Actualiza el estado visual del menú de bandeja
 * 
 * Sincroniza el tooltip y otros elementos visuales con el estado actual
 * de la configuración. Se llama después de cambios de estado.
 * 
 * @private
 */
func (s *SystrayManager) updateMenuState() {
	config := s.controller.GetConfig()
	
	// Actualizar tooltip con información actual
	tooltip := "🌙 Luz Nocturna - " + config.GetTemperatureString()
	if config.IsActive {
		tooltip += " (Activa)"
	} else {
		tooltip += " (Inactiva)"
	}
	
	systray.SetTooltip(tooltip)
}

/**
 * showMainWindow - Muestra y enfoca la ventana principal
 * 
 * Hace visible la ventana principal de la aplicación y la trae al frente.
 * Solo funciona si mainView no es nil.
 * 
 * @private
 */
func (s *SystrayManager) showMainWindow() {
	if s.mainView != nil && s.mainView.window != nil {
		s.mainView.window.Show()
		s.mainView.window.RequestFocus()
		println("📱 Ventana principal mostrada desde bandeja")
	}
}
