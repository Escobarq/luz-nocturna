APP_NAME=luz_nocturna
PACKAGE_ID=com.luznocturna.luz_nocturna

all: build

run:
	go run .

build:
	go build -o bin/$(APP_NAME) .

# Crear iconos desde SVG (requiere ImageMagick)
icon: 
	@if command -v convert >/dev/null 2>&1; then \
		convert internal/views/icons/nightlight_icon.svg -resize 16x16 internal/views/icons/nightlight_icon_16.png; \
		convert internal/views/icons/nightlight_icon.svg -resize 24x24 internal/views/icons/nightlight_icon_24.png; \
		convert internal/views/icons/nightlight_icon.svg -resize 32x32 internal/views/icons/nightlight_icon_32.png; \
		convert internal/views/icons/nightlight_icon.svg -resize 32x32 icon.png; \
		echo "✅ Iconos creados: PNG de 16x16, 24x24, 32x32 e icon.png"; \
	else \
		echo "❌ ImageMagick no disponible. Instala con: sudo apt install imagemagick"; \
		exit 1; \
	fi

# Crear paquete (requiere icono)
package: icon
	fyne package -os linux -icon icon.png -name $(APP_NAME) --app-id $(PACKAGE_ID)
	@echo "✅ Paquete creado: $(APP_NAME).tar.xz"

install: build
	sudo cp bin/$(APP_NAME) /usr/local/bin/$(APP_NAME)
	@echo "✅ Instalado en /usr/local/bin/$(APP_NAME)"

clean:
	rm -rf bin $(APP_NAME) $(APP_NAME).tar.xz icon.png *.deb *.rpm
	@echo "✅ Archivos limpiados"

# Mostrar ayuda
help:
	@echo "🌙 Luz Nocturna - Comandos disponibles:"
	@echo ""
	@echo "  make build    - Compilar binario en bin/"
	@echo "  make run      - Ejecutar en modo desarrollo"  
	@echo "  make icon     - Crear icon.png desde SVG"
	@echo "  make package  - Crear paquete tar.xz con icono"
	@echo "  make install  - Instalar en sistema (requiere sudo)"
	@echo "  make clean    - Limpiar archivos generados"
	@echo "  make help     - Mostrar esta ayuda"

.PHONY: all run build icon package install clean help
