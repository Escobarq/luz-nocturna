package models

import (
	"fmt"
	"time"
)

/**
 * Scheduler - Manejador de programación automática de horarios
 *
 * Maneja la aplicación automática de filtros de luz nocturna basada en
 * horarios configurados por el usuario, con soporte para transiciones
 * suaves entre temperaturas de color.
 */
type Scheduler struct {
	config      *AppConfig
	isRunning   bool
	stopChannel chan bool
	onApply     func(float64) error // Callback para aplicar temperatura
}

/**
 * NewScheduler - Constructor del programador de horarios
 *
 * @param {*AppConfig} config - Configuración de la aplicación
 * @param {func(float64) error} onApply - Función callback para aplicar temperatura
 * @returns {*Scheduler} Nueva instancia del programador
 */
func NewScheduler(config *AppConfig, onApply func(float64) error) *Scheduler {
	return &Scheduler{
		config:      config,
		isRunning:   false,
		stopChannel: make(chan bool),
		onApply:     onApply,
	}
}

/**
 * Start - Inicia el programador automático de horarios
 *
 * Comienza a monitorear la hora actual y aplica automáticamente
 * los filtros de temperatura según la configuración.
 */
func (s *Scheduler) Start() {
	if s.isRunning || !s.config.ScheduleEnabled {
		return
	}

	s.isRunning = true
	fmt.Println("🕐 Programación automática iniciada")

	go func() {
		// Aplicar temperatura inicial inmediatamente
		s.applyCurrentTemperature()

		// Crear ticker para verificar cada minuto
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.applyCurrentTemperature()
			case <-s.stopChannel:
				fmt.Println("🕐 Programación automática detenida")
				return
			}
		}
	}()
}

/**
 * Stop - Detiene el programador automático de horarios
 */
func (s *Scheduler) Stop() {
	if !s.isRunning {
		return
	}

	s.isRunning = false
	s.stopChannel <- true
}

/**
 * IsRunning - Verifica si el programador está ejecutándose
 *
 * @returns {bool} true si está ejecutándose
 */
func (s *Scheduler) IsRunning() bool {
	return s.isRunning
}

/**
 * applyCurrentTemperature - Aplica la temperatura correspondiente a la hora actual
 *
 * Calcula la temperatura que debe aplicarse según la hora actual
 * y los horarios configurados, incluyendo transiciones suaves.
 *
 * @private
 */
func (s *Scheduler) applyCurrentTemperature() {
	now := time.Now()
	currentTime := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

	temperature := s.calculateTemperatureForTime(currentTime)

	if s.onApply != nil {
		if err := s.onApply(temperature); err != nil {
			fmt.Printf("⚠️  Error aplicando temperatura automática: %v\n", err)
		} else {
			fmt.Printf("🕐 Temperatura automática aplicada: %.0fK (%s)\n", temperature, currentTime)
		}
	}
}

/**
 * calculateTemperatureForTime - Calcula la temperatura para una hora específica
 *
 * Determina qué temperatura aplicar basándose en los horarios configurados
 * y aplica transiciones suaves durante los períodos de cambio.
 *
 * @param {string} currentTime - Hora actual en formato "HH:MM"
 * @returns {float64} Temperatura a aplicar en Kelvin
 * @private
 */
func (s *Scheduler) calculateTemperatureForTime(currentTime string) float64 {
	schedule := s.config.Schedule

	// Convertir horarios a minutos desde medianoche para facilitar comparaciones
	currentMinutes := s.timeToMinutes(currentTime)
	startMinutes := s.timeToMinutes(schedule.StartTime)
	endMinutes := s.timeToMinutes(schedule.EndTime)

	// Manejar casos donde el período nocturno cruza medianoche (ej: 20:00 - 07:00)
	var isNightPeriod bool
	if startMinutes > endMinutes {
		// El período nocturno cruza medianoche
		isNightPeriod = currentMinutes >= startMinutes || currentMinutes <= endMinutes
	} else {
		// El período nocturno no cruza medianoche
		isNightPeriod = currentMinutes >= startMinutes && currentMinutes <= endMinutes
	}

	// Calcular si estamos en período de transición
	transitionMinutes := schedule.TransitionTime

	if isNightPeriod {
		// Estamos en período nocturno
		if transitionMinutes > 0 {
			// Verificar si estamos en transición al inicio del período nocturno
			transitionStart := startMinutes
			transitionEnd := (startMinutes + transitionMinutes) % (24 * 60)

			if s.isInTransitionPeriod(currentMinutes, transitionStart, transitionEnd, startMinutes > endMinutes) {
				// Calcular progreso de transición (0.0 = inicio, 1.0 = final)
				progress := s.calculateTransitionProgress(currentMinutes, transitionStart, transitionEnd, startMinutes > endMinutes)
				return s.interpolateTemperature(schedule.DayTemp, schedule.NightTemp, progress)
			}
		}
		return schedule.NightTemp
	} else {
		// Estamos en período diurno
		if transitionMinutes > 0 {
			// Verificar si estamos en transición al final del período nocturno
			transitionStart := (endMinutes - transitionMinutes + 24*60) % (24 * 60)
			transitionEnd := endMinutes

			if s.isInTransitionPeriod(currentMinutes, transitionStart, transitionEnd, startMinutes > endMinutes) {
				// Calcular progreso de transición (0.0 = inicio, 1.0 = final)
				progress := s.calculateTransitionProgress(currentMinutes, transitionStart, transitionEnd, startMinutes > endMinutes)
				return s.interpolateTemperature(schedule.NightTemp, schedule.DayTemp, progress)
			}
		}
		return schedule.DayTemp
	}
}

/**
 * timeToMinutes - Convierte tiempo "HH:MM" a minutos desde medianoche
 *
 * @param {string} timeStr - Tiempo en formato "HH:MM"
 * @returns {int} Minutos desde medianoche
 * @private
 */
func (s *Scheduler) timeToMinutes(timeStr string) int {
	var hours, minutes int
	fmt.Sscanf(timeStr, "%d:%d", &hours, &minutes)
	return hours*60 + minutes
}

/**
 * isInTransitionPeriod - Verifica si estamos en un período de transición
 *
 * @param {int} current - Minutos actuales
 * @param {int} start - Inicio de transición
 * @param {int} end - Final de transición
 * @param {bool} crossesMidnight - Si el período cruza medianoche
 * @returns {bool} true si estamos en transición
 * @private
 */
func (s *Scheduler) isInTransitionPeriod(current, start, end int, crossesMidnight bool) bool {
	if crossesMidnight && start > end {
		return current >= start || current <= end
	}
	return current >= start && current <= end
}

/**
 * calculateTransitionProgress - Calcula el progreso de una transición
 *
 * @param {int} current - Minutos actuales
 * @param {int} start - Inicio de transición
 * @param {int} end - Final de transición
 * @param {bool} crossesMidnight - Si el período cruza medianoche
 * @returns {float64} Progreso de 0.0 a 1.0
 * @private
 */
func (s *Scheduler) calculateTransitionProgress(current, start, end int, crossesMidnight bool) float64 {
	var duration int
	var elapsed int

	if crossesMidnight && start > end {
		duration = (24*60 - start) + end
		if current >= start {
			elapsed = current - start
		} else {
			elapsed = (24*60 - start) + current
		}
	} else {
		duration = end - start
		elapsed = current - start
	}

	if duration <= 0 {
		return 1.0
	}

	progress := float64(elapsed) / float64(duration)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	return progress
}

/**
 * interpolateTemperature - Interpola entre dos temperaturas
 *
 * @param {float64} from - Temperatura inicial
 * @param {float64} to - Temperatura final
 * @param {float64} progress - Progreso (0.0 a 1.0)
 * @returns {float64} Temperatura interpolada
 * @private
 */
func (s *Scheduler) interpolateTemperature(from, to, progress float64) float64 {
	return from + (to-from)*progress
}

/**
 * GetNextScheduleChange - Obtiene información sobre el próximo cambio programado
 *
 * @returns {string, float64, time.Duration} Descripción, temperatura y tiempo restante
 */
func (s *Scheduler) GetNextScheduleChange() (string, float64, time.Duration) {
	if !s.config.ScheduleEnabled {
		return "Programación deshabilitada", s.config.LastTemperature, 0
	}

	now := time.Now()
	schedule := s.config.Schedule

	// Obtener horarios de hoy
	startTime := s.parseTimeToday(schedule.StartTime)
	endTime := s.parseTimeToday(schedule.EndTime)

	// Si el horario de fin es antes que el de inicio, significa que cruza medianoche
	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}

	var nextChange time.Time
	var nextTemp float64
	var description string

	if now.Before(startTime) {
		// Próximo cambio es el inicio del período nocturno
		nextChange = startTime
		nextTemp = schedule.NightTemp
		description = "Inicio filtro nocturno"
	} else if now.Before(endTime) {
		// Estamos en período nocturno, próximo cambio es el fin
		nextChange = endTime
		nextTemp = schedule.DayTemp
		description = "Fin filtro nocturno"
	} else {
		// Próximo cambio es el inicio del día siguiente
		nextChange = startTime.Add(24 * time.Hour)
		nextTemp = schedule.NightTemp
		description = "Inicio filtro nocturno"
	}

	duration := nextChange.Sub(now)
	return description, nextTemp, duration
}

/**
 * parseTimeToday - Convierte "HH:MM" a time.Time para hoy
 *
 * @param {string} timeStr - Tiempo en formato "HH:MM"
 * @returns {time.Time} Tiempo completo para hoy
 * @private
 */
func (s *Scheduler) parseTimeToday(timeStr string) time.Time {
	var hours, minutes int
	fmt.Sscanf(timeStr, "%d:%d", &hours, &minutes)

	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, now.Location())
}

/**
 * UpdateConfig - Actualiza la configuración del programador
 *
 * @param {*AppConfig} newConfig - Nueva configuración
 */
func (s *Scheduler) UpdateConfig(newConfig *AppConfig) {
	s.config = newConfig

	// Si la programación se deshabilitó, detener
	if !newConfig.ScheduleEnabled && s.isRunning {
		s.Stop()
	}

	// Si se habilitó y no está corriendo, iniciar
	if newConfig.ScheduleEnabled && !s.isRunning {
		s.Start()
	}
}
