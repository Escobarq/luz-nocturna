package views

import (
	_ "embed"
)

//go:embed icons/nightlight_icon.svg
var nightlightIconSVG []byte

//go:embed icons/nightlight_icon_16.png
var nightlightIcon16 []byte

//go:embed icons/nightlight_icon_24.png
var nightlightIcon24 []byte

/**
 * GetOptimalIcon - Selecciona el icono más apropiado según el sistema
 *
 * Detecta el entorno de escritorio y devuelve el icono en el formato
 * y tamaño más adecuado para la bandeja del sistema.
 *
 * @returns {[]byte} Datos del icono en el formato óptimo
 */
func GetOptimalIcon() []byte {
	// Preferir PNG para mejor compatibilidad con bandejas del sistema
	// 16x16 es el tamaño más compatible universalmente
	if len(nightlightIcon16) > 0 {
		return nightlightIcon16
	}
	// Fallback a 24x24 si 16x16 no está disponible
	if len(nightlightIcon24) > 0 {
		return nightlightIcon24
	}
	// Último recurso: SVG
	return nightlightIconSVG
}
