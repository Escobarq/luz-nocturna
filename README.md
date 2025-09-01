# 🌙 Luz Nocturna

Una aplicación de escritorio para controlar el filtro de luz nocturna en sistemas Linux, construida con Go y Fyne siguiendo el patrón arquitectural MVC. **Implementación nativa con xrandr, sin dependencias externas**.

## ✨ Características

- ✅ **Control nativo con xrandr** - Sin dependencias de redshift
- ✅ **Interfaz gráfica intuitiva** con Fyne
- ✅ **Control de temperatura de color** (3000K - 6500K)
- ✅ **Presets predefinidos** (Cálida, Neutra, Fría, Diurna)
- ✅ **Bandeja del sistema** con menú contextual
- ✅ **Autostart** y minimizar a bandeja
- ✅ **Detección automática** de displays conectados
- ✅ **Configuración persistente** - Recuerda tus preferencias
- ✅ **Arquitectura MVC** bien organizada
- ✅ **Diálogos auto-cerrables** para mejor UX

## 🚀 Instalación Rápida

### Opción 1: Script de Instalación (Recomendado)
```bash
# Compilar
go build -o luz-nocturna main.go

# Instalar (incluye autostart opcional)
./install.sh
```

### Opción 2: Manual
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
├── install.sh                  # Script de instalación
├── go.mod                      # Dependencias de Go
├── README.md                   # Esta documentación
├── DEVELOPMENT.md              # Guía de desarrollo
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

### 🌡️ Control de Temperatura
- **Slider interactivo**: 3000K (cálida) - 6500K (fría)
- **Presets con un clic**: 🕯️ Cálida, ☀️ Neutra, 🌤️ Fría, ☀️ Diurna
- **Indicador visual** del tipo de temperatura actual
- **Aplicación inmediata** a todos los displays conectados

### 🖥️ Soporte Multi-Display
- **Detección automática** de pantallas conectadas
- **Aplicación simultánea** a todas las pantallas
- **Información visual** de displays detectados

### ⚙️ Configuración Inteligente
- **Persistencia automática** en `~/.config/luz-nocturna/config.json`
- **Recuerda última temperatura** usada
- **Configuración de autostart** y comportamiento

### 🎨 Interfaz de Usuario
- **Diálogos auto-cerrables**: Se cierran automáticamente tras 2 segundos
- **Botón Toggle**: Activar/desactivar rápidamente
- **Información en tiempo real**: Estado y displays conectados
- **Diseño responsive**: Se adapta al contenido

## 🔧 Implementación Técnica

### Sistema Gamma Nativo
La aplicación usa `xrandr` directamente para controlar la temperatura de color:

```bash
# Ejemplo de comando generado internamente
xrandr --output eDP-1 --gamma 1.0:0.8:0.6
```

### Algoritmo de Conversión
- **Conversión Kelvin → RGB** usando algoritmo optimizado
- **Rangos seguros** para evitar valores extremos
- **Aplicación por display** individual

### Arquitectura MVC
- **Modelos**: Lógica de negocio y persistencia
- **Vistas**: UI con Fyne + bandeja del sistema
- **Controladores**: Coordinación entre modelo y vista

## 🛠️ Dependencias

### Sistema
- **Linux** con X11 (requerido para xrandr)
- **xrandr** (usualmente incluido)
- **Entorno de escritorio** con soporte para bandeja del sistema

### Go Módulos
- **fyne.io/fyne/v2** - Framework UI
- **fyne.io/systray** - Soporte bandeja del sistema
- **Go 1.22+** - Lenguaje base

### Verificar Sistema
```bash
# Verificar xrandr
xrandr --version

# Ver displays disponibles
xrandr | grep connected
```

## 🔮 Próximas Mejoras

- 🕐 **Programación automática** por horario
- 🌍 **Detección de ubicación** para sunrise/sunset
- 📊 **Perfiles personalizados** con nombres propios
- 🎨 **Temas visuales** y personalización
- 📦 **Paquetes .deb/.rpm** para distribución
- 🔄 **Actualizaciones automáticas**

## 🐛 Solución de Problemas

### La temperatura no se aplica
```bash
# Verificar xrandr funciona
xrandr --output eDP-1 --gamma 1.0:0.8:0.6

# Ver displays disponibles
xrandr | grep connected
```

### No aparece en bandeja del sistema
- Verifica que tu escritorio soporte bandejas del sistema
- En GNOME: instala extensión "AppIndicator Support"
- En KDE/XFCE: Soporte nativo

### Problemas de permisos
```bash
# Asegurar permisos correctos
chmod +x /usr/local/bin/luz-nocturna
```

## 📄 Licencia

MIT - Libre para uso personal y comercial

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! 
- 🐛 Reporta bugs
- 💡 Sugiere mejoras  
- 🔧 Envía pull requests

---
**💡 Tip**: Usa `luz-nocturna --tray` para ejecutar discretamente en segundo plano.
