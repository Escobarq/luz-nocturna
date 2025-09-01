package styles

import (
	"image/color"
)

// Colores principales de la aplicación
var (
	// Colores de fondo
	BackgroundColor = color.NRGBA{R: 242, G: 242, B: 242, A: 255} // #f2f2f2
	WindowColor     = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // #fff

	// Colores de texto
	PrimaryTextColor   = color.NRGBA{R: 51, G: 51, B: 51, A: 255}    // #333
	SecondaryTextColor = color.NRGBA{R: 102, G: 102, B: 102, A: 255} // #666
	TitleTextColor     = color.NRGBA{R: 0, G: 0, B: 0, A: 255}       // #000

	// Colores de botones
	PrimaryButtonColor        = color.NRGBA{R: 0, G: 120, B: 212, A: 255}   // #0078d4
	PrimaryButtonHoverColor   = color.NRGBA{R: 0, G: 90, B: 158, A: 255}    // #005a9e
	SecondaryButtonColor      = color.NRGBA{R: 221, G: 221, B: 221, A: 255} // #ddd
	SecondaryButtonHoverColor = color.NRGBA{R: 187, G: 187, B: 187, A: 255} // #bbb

	// Colores de slider
	SliderBackgroundColor = color.NRGBA{R: 230, G: 230, B: 230, A: 255} // #e6e6e6
	SliderActiveColor     = color.NRGBA{R: 0, G: 120, B: 212, A: 255}   // #0078d4

	// Sombras
	ShadowColor = color.NRGBA{R: 0, G: 0, B: 0, A: 25} // rgba(0,0,0,0.1)
)
