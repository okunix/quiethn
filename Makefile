GO_MAIN=main.go
GO_SOURCES=$(shell find . -name '*.go' -not -path './vendor/*')

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
	CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -o $@ $<

.PHONY: run 
run: ${BIN_PATH}
	if [ ! -x $< ]; then chmod +x $<; fi
	HN_SERVER_PORT=8080 $<

.PHONY: clean
clean:
	go clean
	rm -rf ./bin

DOCKER_IMAGE=ghcr.io/okunix/quiethn
DOCKER_TAG=latest
DOCKERFILE=Dockerfile

.PHONY: docker-build
docker-build: ${DOCKERFILE} ${GO_SOURCES}
	docker buildx build -t ${DOCKER_IMAGE}:${DOCKER_TAG} -f ${DOCKERFILE} .

.PHONY: docker-push
docker-push: docker-build 
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}

.PHONY: docker-run
docker-run: docker-build
	docker run --rm -p 8080:80 ${DOCKER_IMAGE}:${DOCKER_TAG}

.PHONY: docker-clean
docker-clean:
	docker rmi ${DOCKER_IMAGE}:${DOCKER_TAG}
