# 🛠️ Guía de Desarrollo - Luz Nocturna

## Instalación de Dependencias del Sistema

Para que la aplicación funcione completamente, necesita tener instalado `redshift`:

```bash
# Ubuntu/Debian
sudo apt install redshift

# Fedora
sudo dnf install redshift

# Arch Linux
sudo pacman -S redshift
```

## Estructura MVC Implementada

### 📊 Modelos (Models)
- **`nightlight.go`**: Lógica de negocio principal
- **`config.go`**: Configuración persistente de la aplicación
- **`presets.go`**: Presets de temperatura predefinidos

### 🎨 Vistas (Views)
- **`nightlight_view.go`**: Interfaz gráfica principal con Fyne

### 🎮 Controladores (Controllers)
- **`nightlight_controller.go`**: Coordina modelo y vista

### 🎨 Estilos (Styles)
- **`colors.go`**: Paleta de colores
- **`dimensions.go`**: Dimensiones y estilos

### ⚙️ Sistema (System)
- **`gamma.go`**: Integración con herramientas de gamma del sistema

## Ejemplos de Uso del Código

### Usar el modelo directamente
```go
config := models.NewNightLightConfig()
config.SetTemperature(3500)
fmt.Println(config.GetTemperatureString()) // "3500K"
```

### Usar presets de temperatura
```go
temp := models.Presets.GetRecommendedForTime(22) // 10 PM
fmt.Println(temp) // 3000 (temperatura cálida para la noche)
```

### Aplicar configuración del sistema
```go
gm := system.NewGammaManager()
if gm.IsRedshiftInstalled() {
    gm.ApplyTemperature(4000)
}
```

## Mejoras Futuras Implementables

### 1. Modo Automático por Horario
```go
// En models/scheduler.go
type AutoScheduler struct {
    config *NightLightConfig
    enabled bool
}

func (s *AutoScheduler) UpdateByTime() {
    now := time.Now()
    recommended := Presets.GetRecommendedForTime(now.Hour())
    s.config.SetTemperature(recommended)
}
```

### 2. Icono de Bandeja del Sistema
```go
// En views/systray.go usando fyne.io/systray
func (v *NightLightView) setupSystray() {
    systray.Run(v.onSystrayReady, v.onSystrayExit)
}
```

### 3. Configuración Avanzada
```go
// En models/advanced_config.go
type AdvancedConfig struct {
    SunriseTime    time.Time
    SunsetTime     time.Time
    AutoTransition bool
    TransitionDuration time.Duration
}
```

### 4. Detección Automática de Ubicación
```go
// En system/location.go
type LocationService struct{}

func (l *LocationService) GetSunriseSunset() (sunrise, sunset time.Time, err error) {
    // Integración con API de sunrise-sunset o geolocalización
}
```

## Comandos de Desarrollo

### Ejecutar en modo desarrollo con logs
```bash
go run main.go 2>&1 | tee debug.log
```

### Compilar para diferentes arquitecturas
```bash
# Linux x64
GOOS=linux GOARCH=amd64 go build -o luz-nocturna-linux-x64 main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o luz-nocturna-linux-arm64 main.go
```

### Generar paquete .deb (requiere fpm)
```bash
# Instalar fpm
gem install fpm

# Crear paquete
fpm -s dir -t deb -n luz-nocturna -v 1.0.0 \
    --description "Aplicación de luz nocturna para Linux" \
    --license MIT \
    --vendor "Tu Nombre" \
    luz-nocturna=/usr/local/bin/luz-nocturna
```

## Testing

### Estructura de tests recomendada
```
internal/
├── controllers/
│   ├── nightlight_controller.go
│   └── nightlight_controller_test.go
├── models/
│   ├── nightlight.go
│   └── nightlight_test.go
```

### Ejemplo de test unitario
```go
// En internal/models/nightlight_test.go
func TestNightLightConfig_SetTemperature(t *testing.T) {
    config := NewNightLightConfig()
    
    config.SetTemperature(5000)
    assert.Equal(t, 5000.0, config.Temperature)
    
    // Probar límites
    config.SetTemperature(2000) // Muy bajo
    assert.Equal(t, 3000.0, config.Temperature) // Debe ser el mínimo
}
```

## Integración con Escritorio

### Archivo .desktop
```ini
[Desktop Entry]
Version=1.0
Type=Application
Name=Luz Nocturna
Comment=Control de temperatura de color
Exec=/usr/local/bin/luz-nocturna
Icon=luz-nocturna
Terminal=false
StartupNotify=true
Categories=Utility;System;
```

### Autostart
```bash
# Crear entrada de autostart
mkdir -p ~/.config/autostart
cp luz-nocturna.desktop ~/.config/autostart/
```

## Consideraciones de Performance

- Los cambios de gamma son operaciones del sistema que pueden ser costosas
- Implementar debouncing para el slider (evitar aplicar cada cambio inmediatamente)
- Cache de configuración para evitar I/O excesivo
- Lazy loading de recursos pesados
