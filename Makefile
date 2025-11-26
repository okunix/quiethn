GO_MAIN=main.go
GO_SOURCES=$(shell find . -name *.go -not -path './vendor/*')

TARGET_ARCH=amd64
TARGET_OS=linux

EXECUTABLE_NAME=quiethn
BIN_PATH=./bin/${EXECUTABLE_NAME}

.PHONY: all
all: build

.PHONY: build
build: ${BIN_PATH}

${BIN_PATH}: ${GO_MAIN} ${GO_SOURCES}
	go mod tidy
	go mod download
	go test ./...
	GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -o $@ $<
	if [ -f config.yaml ]; then cp config.yaml ./bin/config.yaml; fi 

.PHONY: run 
run: ${BIN_PATH}
	chmod +x $<
	$<

.PHONY: clean
clean:
	go clean
	rm -rf ./bin
