PACKAGE := github.com/axelrindle/proxyguy
TIMESTAMP := $(shell date +'%Y-%m-%d %T')
VERSION ?= nightly

OUTPUT_DIR := dist
OUTPUT_FILE := proxyguy

OUTPUT ?= $(OUTPUT_DIR)/$(OUTPUT_FILE)

default: build

clean:
	@rm -rf $(OUTPUT_DIR)

build: clean
	go build -v \
		-ldflags="-s -w -X '$(PACKAGE)/cli.Version=$(VERSION)' -X '$(PACKAGE)/cli.BuildTime=$(TIMESTAMP)'" \
		-o $(OUTPUT) .

build-static:
	CGO_ENABLED=0 make build

build-all: build build-static

test:
	go test -v -cover -coverprofile=coverage.out ./...

run: build
	./dist/proxyguy

run-local:
	go run . --config ./config.local.yml

install:
	install dist/proxyguy /usr/bin
	if [ ! -d /etc/proxyguy ]; then mkdir /etc/proxyguy; fi
	if [ ! -f /etc/proxyguy/config.yaml ]; then install -m 644 -b -T config.example.yml /etc/proxyguy/config.yaml; fi
