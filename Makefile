OUTPUT ?= dist


default: build

clean:
	@rm -rf $(OUTPUT)

build: clean
	go build -ldflags="-s -w" -o $(OUTPUT)/proxyguy .

test:
	go test -v -cover -coverprofile=coverage.out ./...

run: build
	./dist/proxyguy

run-local:
	go run . -config ./config.test.yml

install: build
	if [ ! -d /etc/proxyguy ]; then mkdir /etc/proxyguy; fi
	install dist/proxyguy /usr/bin
	if [ ! -f /etc/proxyguy/config.yaml ]; then install -m 644 -b -T config.example.yml /etc/proxyguy/config.yaml; fi
