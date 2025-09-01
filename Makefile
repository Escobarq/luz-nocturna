APP_NAME=luz_nocturna
DEB_NAME=luz-nocturna
PACKAGE_ID=com.luznocturna.luz_nocturna
VERSION=1.0.0
MAINTAINER=Juan <juan@example.com>
DESCRIPTION=Control de temperatura de color para monitores Linux
HOMEPAGE=https://github.com/juan/luz-nocturna

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

# Crear paquete Fyne (requiere icono)
package: icon
	fyne package -os linux -icon icon.png -name $(APP_NAME) --app-id $(PACKAGE_ID)
	@echo "✅ Paquete Fyne creado: $(APP_NAME).tar.xz"

# Crear paquete DEB para Debian/Ubuntu
deb: build icon
	@echo "📦 Creando paquete DEB..."
	@rm -rf deb_build
	@mkdir -p deb_build/DEBIAN
	@mkdir -p deb_build/usr/local/bin
	@mkdir -p deb_build/usr/share/applications
	@mkdir -p deb_build/usr/share/pixmaps
	@mkdir -p deb_build/usr/share/doc/$(DEB_NAME)
	
	# Copiar binario
	@cp bin/$(APP_NAME) deb_build/usr/local/bin/
	@chmod 755 deb_build/usr/local/bin/$(APP_NAME)
	
	# Copiar icono
	@cp icon.png deb_build/usr/share/pixmaps/$(DEB_NAME).png
	
	# Crear archivo .desktop
	@echo "[Desktop Entry]" > deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "Name=Luz Nocturna" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "Comment=$(DESCRIPTION)" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "Exec=/usr/local/bin/$(APP_NAME)" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "Icon=/usr/share/pixmaps/$(DEB_NAME).png" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "Terminal=false" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "Type=Application" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "Categories=Utility;Settings;" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@echo "StartupNotify=true" >> deb_build/usr/share/applications/$(DEB_NAME).desktop
	@chmod 644 deb_build/usr/share/applications/$(DEB_NAME).desktop
	
	# Crear documentación
	@echo "$(DESCRIPTION)" > deb_build/usr/share/doc/$(DEB_NAME)/README
	@echo "" >> deb_build/usr/share/doc/$(DEB_NAME)/README
	@echo "Aplicación desarrollada en Go con Fyne.io" >> deb_build/usr/share/doc/$(DEB_NAME)/README
	@echo "Versión: $(VERSION)" >> deb_build/usr/share/doc/$(DEB_NAME)/README
	@gzip -9c deb_build/usr/share/doc/$(DEB_NAME)/README > deb_build/usr/share/doc/$(DEB_NAME)/README.gz
	@rm deb_build/usr/share/doc/$(DEB_NAME)/README
	
	# Crear archivo de control
	@echo "Package: $(DEB_NAME)" > deb_build/DEBIAN/control
	@echo "Version: $(VERSION)" >> deb_build/DEBIAN/control
	@echo "Section: utils" >> deb_build/DEBIAN/control
	@echo "Priority: optional" >> deb_build/DEBIAN/control
	@echo "Architecture: amd64" >> deb_build/DEBIAN/control
	@echo "Depends: libc6 (>= 2.31)" >> deb_build/DEBIAN/control
	@echo "Maintainer: $(MAINTAINER)" >> deb_build/DEBIAN/control
	@echo "Description: $(DESCRIPTION)" >> deb_build/DEBIAN/control
	@echo " Una aplicación moderna para controlar la temperatura de color" >> deb_build/DEBIAN/control
	@echo " de los monitores en sistemas Linux. Desarrollada con Go y Fyne." >> deb_build/DEBIAN/control
	@echo " ." >> deb_build/DEBIAN/control  
	@echo " Características:" >> deb_build/DEBIAN/control
	@echo "  * Interfaz gráfica intuitiva" >> deb_build/DEBIAN/control
	@echo "  * Soporte para bandeja del sistema" >> deb_build/DEBIAN/control
	@echo "  * Presets de temperatura predefinidos" >> deb_build/DEBIAN/control
	@echo "  * Compatible con X11 y Wayland" >> deb_build/DEBIAN/control
	@echo "Homepage: $(HOMEPAGE)" >> deb_build/DEBIAN/control
	
	# Construir el paquete
	@dpkg-deb --build deb_build $(DEB_NAME)_$(VERSION)_amd64.deb
	@rm -rf deb_build
	@echo "✅ Paquete DEB creado: $(DEB_NAME)_$(VERSION)_amd64.deb"

# Crear paquete RPM para RedHat/Fedora
rpm: build icon
	@echo "📦 Creando paquete RPM..."
	@if ! command -v rpmbuild >/dev/null 2>&1; then \
		echo "❌ rpmbuild no disponible. Instala con:"; \
		echo "   Fedora: sudo dnf install rpm-build"; \
		echo "   Ubuntu: sudo apt install rpm"; \
		exit 1; \
	fi
	@rm -rf rpm_build ~/rpmbuild
	@mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
	@mkdir -p rpm_build/$(DEB_NAME)-$(VERSION)/usr/local/bin
	@mkdir -p rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications
	@mkdir -p rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/pixmaps
	
	# Copiar archivos
	@cp bin/$(APP_NAME) rpm_build/$(DEB_NAME)-$(VERSION)/usr/local/bin/
	@cp icon.png rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/pixmaps/$(DEB_NAME).png
	@cp deb_build_temp/usr/share/applications/$(DEB_NAME).desktop rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/ || echo "Creando .desktop para RPM..."
	
	# Recrear .desktop para RPM
	@echo "[Desktop Entry]" > rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	@echo "Name=Luz Nocturna" >> rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	@echo "Comment=$(DESCRIPTION)" >> rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	@echo "Exec=/usr/local/bin/$(APP_NAME)" >> rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	@echo "Icon=/usr/share/pixmaps/$(DEB_NAME).png" >> rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	@echo "Terminal=false" >> rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	@echo "Type=Application" >> rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	@echo "Categories=Utility;Settings;" >> rpm_build/$(DEB_NAME)-$(VERSION)/usr/share/applications/$(DEB_NAME).desktop
	
	# Crear tarball fuente
	@cd rpm_build && tar -czf ~/rpmbuild/SOURCES/$(DEB_NAME)-$(VERSION).tar.gz $(DEB_NAME)-$(VERSION)/
	
	# Crear spec file
	@echo "Name: $(DEB_NAME)" > ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "Version: $(VERSION)" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "Release: 1%{?dist}" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "Summary: $(DESCRIPTION)" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "License: MIT" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "URL: $(HOMEPAGE)" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "Source0: %{name}-%{version}.tar.gz" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "BuildArch: x86_64" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "Requires: glibc" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "%description" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "Una aplicación moderna para controlar la temperatura de color" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "de los monitores en sistemas Linux. Desarrollada con Go y Fyne." >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "%prep" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "%setup -q" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "%install" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "mkdir -p %{buildroot}/usr/local/bin" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "mkdir -p %{buildroot}/usr/share/applications" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "mkdir -p %{buildroot}/usr/share/pixmaps" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "cp -p usr/local/bin/$(APP_NAME) %{buildroot}/usr/local/bin/" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "cp -p usr/share/applications/$(DEB_NAME).desktop %{buildroot}/usr/share/applications/" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "cp -p usr/share/pixmaps/$(DEB_NAME).png %{buildroot}/usr/share/pixmaps/" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "%files" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "/usr/local/bin/$(APP_NAME)" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "/usr/share/applications/$(DEB_NAME).desktop" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@echo "/usr/share/pixmaps/$(DEB_NAME).png" >> ~/rpmbuild/SPECS/$(DEB_NAME).spec
	
	# Construir RPM
	@rpmbuild -ba ~/rpmbuild/SPECS/$(DEB_NAME).spec
	@cp ~/rpmbuild/RPMS/x86_64/$(DEB_NAME)-$(VERSION)-1.*.x86_64.rpm ./
	@rm -rf rpm_build ~/rpmbuild
	@echo "✅ Paquete RPM creado: $(DEB_NAME)-$(VERSION)-1.*.x86_64.rpm"

# Crear todos los paquetes
packages: package deb
	@if command -v rpmbuild >/dev/null 2>&1; then \
		make rpm; \
		echo "✅ Todos los paquetes creados: Fyne (.tar.xz), DEB y RPM"; \
	else \
		echo "✅ Paquetes Fyne (.tar.xz) y DEB creados"; \
		echo "💡 Para crear RPM instala: sudo apt install rpm o sudo dnf install rpm-build"; \
	fi

install: build
	sudo cp bin/$(APP_NAME) /usr/local/bin/$(APP_NAME)
	@echo "✅ Instalado en /usr/local/bin/$(APP_NAME)"

clean:
	rm -rf bin $(APP_NAME) $(APP_NAME).tar.xz icon.png *.deb *.rpm deb_build rpm_build
	@echo "✅ Archivos limpiados"

# Mostrar ayuda
help:
	@echo "🌙 Luz Nocturna - Comandos disponibles:"
	@echo ""
	@echo "  make build    - Compilar binario en bin/"
	@echo "  make run      - Ejecutar en modo desarrollo"  
	@echo "  make icon     - Crear icon.png desde SVG"
	@echo "  make package  - Crear paquete Fyne tar.xz con icono"
	@echo "  make deb      - Crear paquete DEB para Debian/Ubuntu" 
	@echo "  make rpm      - Crear paquete RPM para RedHat/Fedora (requiere rpm-build)"
	@echo "  make packages - Crear todos los paquetes (Fyne, DEB, RPM)"
	@echo "  make install  - Instalar en sistema (requiere sudo)"
	@echo "  make clean    - Limpiar archivos generados"
	@echo "  make help     - Mostrar esta ayuda"

.PHONY: all run build icon package deb rpm packages install clean help
