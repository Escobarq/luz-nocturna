# 🌙 Luz Nocturna

Una aplicación de escritorio para controlar el filtro de luz nocturna en sistemas Linux, construida con Go y Fyne siguiendo el patrón arquitectural MVC (Modelo-Vista-Controlador).

## Estructura del Proyecto

```
luz-nocturna/
├── main.go                     # Punto de entrada de la aplicación
├── go.mod                      # Dependencias de Go
├── go.sum                      # Checksums de dependencias
├── icon.png                    # Icono de la aplicación
├── index.html                  # Prototipo HTML original
├── Makefile                    # Comandos de compilación
└── internal/                   # Código interno de la aplicación
    ├── controllers/            # Lógica de control (MVC)
    │   └── nightlight_controller.go
    ├── models/                 # Modelos de datos (MVC)
    │   └── nightlight.go
    ├── styles/                 # Estilos y constantes de diseño
    │   ├── colors.go          # Colores de la aplicación
    │   └── dimensions.go      # Dimensiones y estilos
    └── views/                  # Vistas/UI (MVC)
        └── nightlight_view.go
```

## Arquitectura MVC

### Modelos (`internal/models/`)
- **`nightlight.go`**: Define la estructura de datos `NightLightConfig` que maneja:
  - Temperatura de color (3000K - 6500K)
  - Estado activo/inactivo
  - Métodos para aplicar/resetear configuración

### Vistas (`internal/views/`)
- **`nightlight_view.go`**: Interfaz gráfica con Fyne que incluye:
  - Título con emoji 🌙
  - Slider para ajustar temperatura de color
  - Etiqueta que muestra temperatura actual
  - Botones "Aplicar" y "Reset"

### Controladores (`internal/controllers/`)
- **`nightlight_controller.go`**: Lógica de negocio que:
  - Coordina entre modelo y vista
  - Maneja eventos de la interfaz
  - Aplica validaciones de datos

### Estilos (`internal/styles/`)
- **`colors.go`**: Define la paleta de colores de la aplicación
- **`dimensions.go`**: Constantes de tamaños, padding y estilos

## Compilación y Ejecución

### Ejecutar en modo desarrollo
```bash
go run main.go
```

### Compilar para producción
```bash
go build -o luz-nocturna main.go
```

### Usando Makefile (si existe)
```bash
make build    # Compilar
make run      # Ejecutar
make clean    # Limpiar
```

## Características

- ✅ Interfaz gráfica intuitiva con Fyne
- ✅ Control de temperatura de color (3000K - 6500K)
- ✅ Arquitectura MVC bien organizada
- ✅ Separación clara de responsabilidades
- ✅ Estilos centralizados y reutilizables
- 🔄 Integración con sistema gamma (pendiente)

## Dependencias

- **Fyne v2.6.3**: Framework para interfaces gráficas multiplataforma
- **Go 1.22.2+**: Lenguaje de programación

## Próximas Mejoras

1. **Integración con sistema gamma**: Implementar la aplicación real del filtro de color
2. **Persistencia de configuración**: Guardar preferencias del usuario
3. **Inicio automático**: Opción para ejecutar al inicio del sistema
4. **Programación horaria**: Activar automáticamente según la hora
5. **Icono de bandeja**: Ejecutar en background con icono en la bandeja del sistema

## Desarrollo

El proyecto sigue las mejores prácticas de Go:
- Código organizado en paquetes internos
- Separación clara de responsabilidades
- Documentación en código
- Estructura modular y extensible
