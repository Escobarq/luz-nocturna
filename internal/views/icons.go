package views

import (
	_ "embed"
)

/**
 * ICONOS DEL SISTEMA
 *
 * Los iconos para la bandeja del sistema suelen usar estos tamaños:
 * - 16x16px: Tamaño estándar para la mayoría de bandejas del sistema
 * - 22x22px: Tamaño alternativo en algunos entornos
 * - 24x24px: Tamaño para iconos de alta resolución
 * - 32x32px: Para pantallas de alta densidad (HiDPI)
 *
 * El formato SVG es ideal porque se escala automáticamente a cualquier tamaño.
 * También puedes usar PNG en múltiples resoluciones (16x16, 24x24, 32x32).
 */

//go:embed icons/nightlight_icon.svg
var nightlightIconSVG []byte

//go:embed icons/nightlight_icon_16.png
var nightlightIcon16 []byte

//go:embed icons/nightlight_icon_24.png
var nightlightIcon24 []byte

//go:embed icons/nightlight_icon_32.png
var nightlightIcon32 []byte

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

/**
 * GetIconBySize - Obtiene icono por tamaño específico
 *
 * @param {int} size - Tamaño deseado en píxeles (16, 24, 32)
 * @returns {[]byte} Datos del icono en el tamaño solicitado
 */
func GetIconBySize(size int) []byte {
	switch size {
	case 16:
		return nightlightIcon16
	case 24:
		return nightlightIcon24
	case 32:
		return nightlightIcon32
	default:
		return nightlightIconSVG // SVG se escala a cualquier tamaño
	}
}
