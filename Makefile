APP_NAME=luz_nocturna
PACKAGE_ID=com.luznocturna.luz_nocturna

all: run

run:
	go run .

build:
	go build -o bin/$(APP_NAME) .

package:
	fyne package -os linux -icon icon.png -name $(APP_NAME) --app-id $(PACKAGE_ID)

install:
	sudo cp bin/$(APP_NAME) /usr/local/bin/$(APP_NAME)

clean:
	rm -rf bin $(APP_NAME) *.tar.xz *.deb *.rpm
