# 🌙 Luz Nocturna

Una aplicación de escritorio para controlar el filtro de luz nocturna en sistemas Linux, construida con Go y Fyne siguiendo el patrón arquitectural MVC. **Implementación nativa con xrandr, sin dependencias externas**.

## ✨ Características

- ✅ **Control nativo con xrandr/Wayland** - Soporte completo X11 y Wayland
- ✅ **Optimizado para ZorinOS** - Deshabilita automáticamente el sistema nativo
- ✅ **Interfaz gráfica intuitiva** con Fyne
- ✅ **Control de temperatura de color** (3000K - 6500K)
- ✅ **Presets predefinidos** (Cálida, Neutra, Fría, Diurna)
- ✅ **Bandeja del sistema** con menú contextual
- ✅ **Programación automática por horario** - Transiciones suaves día/noche
- ✅ **Control exclusivo** - Evita conflictos con sistemas nativos
- ✅ **Detección automática** de displays conectados
- ✅ **Configuración persistente** - Recuerda tus preferencias
- ✅ **Arquitectura MVC** bien organizada
- ✅ **Instalación automática** de dependencias Wayland

## 🚀 Instalación Rápida

### Manual
```bash
# Compilar
go build -o luz-nocturna main.go

# Copiar a directorio del sistema
sudo cp luz-nocturna /usr/local/bin/

# Ejecutar
luz-nocturna
```

## 📋 Uso

### Interfaz Gráfica
```bash
luz-nocturna                    # Ventana principal
```

### Solo Bandeja del Sistema
```bash
luz-nocturna --tray            # Solo icono en bandeja
```

### Desde la Bandeja
- **Clic derecho** en el icono para acceder al menú
- **Presets rápidos**: Cálida, Neutra, Fría, Diurna
- **Acciones**: Aplicar, Reset, Mostrar ventana
- **Control de temperatura** sin abrir ventana

## 🏗️ Estructura del Proyecto

```
luz-nocturna/
├── main.go                     # Punto de entrada
├── go.mod                      # Dependencias de Go
├── README.md                   # Esta documentación
└── internal/                   # Código interno
    ├── controllers/            # 🎮 Controladores (MVC)
    │   └── nightlight_controller.go
    ├── models/                 # 📊 Modelos (MVC)
    │   ├── nightlight.go       # Lógica principal
    │   ├── config.go           # Configuración persistente
    │   └── presets.go          # Presets de temperatura
    ├── styles/                 # 🎨 Estilos y colores
    │   ├── colors.go
    │   └── dimensions.go
    ├── system/                 # ⚙️ Integración sistema
    │   └── gamma.go            # Control xrandr nativo
    └── views/                  # 🖼️ Vistas (MVC)
        ├── nightlight_view.go  # UI principal
        └── systray.go          # Bandeja del sistema
```

## 🎯 Funcionalidades Detalladas

### 🕐 Programación Automática por Horario
- **Horarios personalizables**: Define inicio y fin del filtro nocturno
- **Temperaturas independientes**: Configura temperatura diurna (ej: 6500K) y nocturna (ej: 3200K)
- **Transiciones suaves**: Cambios graduales entre temperaturas (0-60 minutos)
- **Aplicación automática**: Se ejecuta en segundo plano sin intervención
- **Información en tiempo real**: Próximo cambio programado y tiempo restante
- **Períodos que cruzan medianoche**: Soporte completo para horarios como 20:00 - 07:00

### 🌡️ Control Manual de Temperatura
- **Slider interactivo**: 3000K (cálida) - 6500K (fría)
- **Presets con un clic**: 🕯️ Cálida, ☀️ Neutra, 🌤️ Fría, ☀️ Diurna
- **Override automático**: Control manual temporal sobre programación automática

### 🖥️ Soporte Multi-Plataforma
- **X11 con xrandr**: Soporte nativo y optimizado
- **Wayland completo**: wl-gamma-relay, wlsunset, gammastep
- **Instalación automática**: Detecta distribución e instala dependencias
- **Detección automática** de displays y protocolo

### ⚙️ Configuración Persistente
- **Archivo de configuración**: `~/.config/luz-nocturna/config.json`
- **Programación guardada**: Horarios y temperaturas se mantienen entre sesiones
- **Autostart opcional**: Iniciar con el sistema y programación automática

## 🔧 Implementación Técnica

### Sistema de Programación Automática
La aplicación incluye un scheduler avanzado que:

```go
// Ejemplo de configuración automática
scheduler := models.NewScheduler(config, gammaManager.ApplyTemperature)
scheduler.Start() // Inicia programación automática

// Calcula temperatura según hora actual
temp := scheduler.CalculateTemperatureForTime("22:30")
// Resultado: transición suave hacia temperatura nocturna
```

### Algoritmo de Transición
- **Interpolación lineal** entre temperaturas día/noche
- **Cálculo de períodos**: Manejo correcto de horarios que cruzan medianoche
- **Verificación por minuto**: Precisión temporal sin consumo excesivo de recursos
- **Progreso de transición**: 0.0 (inicio) a 1.0 (final) para cambios suaves

### Soporte Wayland Mejorado
- **Detección automática** de herramientas disponibles
- **Instalación asistida** con pkexec para permisos
- **Múltiples backends**: wl-gamma-relay, wlsunset, gammastep
- **Fallbacks inteligentes**: Si una herramienta falla, prueba la siguiente

## 🛠️ Dependencias

### Sistema
- **Linux** con X11 o Wayland
- **Para X11**: xrandr (usualmente incluido)
- **Para Wayland**: Una de estas herramientas (se instala automáticamente):
  - `wl-gamma-relay`
  - `wlsunset` 
  - `gammastep`

### Go Módulos
- **fyne.io/fyne/v2** - Framework UI
- **fyne.io/systray** - Soporte bandeja del sistema
- **Go 1.22+** - Lenguaje base

### Verificar Sistema
```bash
# Verificar protocolo en uso
echo $XDG_SESSION_TYPE

# Para X11 - verificar xrandr
xrandr --version && xrandr | grep connected

# Para Wayland - verificar herramientas (se instalan automáticamente)
which wlsunset || which gammastep || which wl-gamma-relay
```

## ⚙️ Configuración de Programación Automática

### Configuración Básica
1. **Abrir la aplicación**: `luz-nocturna`
2. **Habilitar programación**: Marcar checkbox "🕐 Programación automática"
3. **Configurar horarios**:
   - **Inicio**: Hora de activación del filtro nocturno (ej: "20:00")
   - **Fin**: Hora de desactivación del filtro nocturno (ej: "07:00")
4. **Ajustar temperaturas**:
   - **Nocturna**: Temperatura cálida para la noche (ej: 3200K)
   - **Diurna**: Temperatura fría para el día (ej: 6500K)
5. **Tiempo de transición**: Duración del cambio gradual (ej: 30 minutos)

### Ejemplo de Configuración
```json
{
  "schedule_enabled": true,
  "schedule": {
    "start_time": "20:00",
    "end_time": "07:00", 
    "night_temp": 3200,
    "day_temp": 6500,
    "transition_time": 30
  }
}
```

### Comportamiento Automático
- **20:00**: Inicio de transición gradual hacia 3200K (30 minutos)
- **20:30**: Temperatura nocturna completa (3200K)
- **06:30**: Inicio de transición gradual hacia 6500K (30 minutos)  
- **07:00**: Temperatura diurna completa (6500K)

## 🐛 Solución de Problemas

### Error en Wayland: "no se pudo aplicar gamma"
```bash
# Instalar dependencias manualmente si la instalación automática falla
# Para ZorinOS/Ubuntu:
sudo apt install wlsunset

# Para Fedora:
sudo dnf install wlsunset

# Para Arch:
sudo pacman -S wlsunset

# Verificar instalación
which wlsunset
```

### La programación automática no funciona
- Verificar que esté habilitada en la interfaz
- Revisar formato de horarios (debe ser "HH:MM")
- Comprobar que los horarios sean válidos (00:00 - 23:59)
- Verificar archivo de configuración: `~/.config/luz-nocturna/config.json`

### La temperatura no se aplica en X11
```bash
# Verificar xrandr funciona
xrandr --output eDP-1 --gamma 1.0:0.8:0.6

# Ver displays disponibles
xrandr | grep connected
```

### No aparece en bandeja del sistema
- En GNOME: instala extensión "AppIndicator Support"
- En KDE/XFCE: Soporte nativo
- Verificar que el escritorio soporte bandejas del sistema

## 📄 Licencia

MIT - Libre para uso personal y comercial

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! 
- 🐛 Reporta bugs
- 💡 Sugiere mejoras  
- 🔧 Envía pull requests

---
**💡 Tips**: 
- Usa `luz-nocturna --tray` para ejecutar solo en la bandeja del sistema
- La programación automática funciona en segundo plano incluso con la ventana cerrada
- Los cambios de configuración se aplican inmediatamente sin reiniciar
