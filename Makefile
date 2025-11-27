BIN_DIR=./bin

GO_MAIN=cmd/quiethn/main.go
GO_SOURCES=$(shell find . -name '*.go' -not -path './vendor/*')
TARGET_ARCH=amd64
TARGET_OS=linux
EXECUTABLE_NAME=quiethn
BIN=${BIN_DIR}/${EXECUTABLE_NAME}

UPDATER_GO_MAIN=./cmd/quiethn-updater/main.go
UPDATER_GO_SOURCES=${GO_SOURCES}
UPDATER_TARGET_OS=${TARGET_OS}
UPDATER_TARGET_ARCH=${TARGET_ARCH}
UPDATER_EXECUTABLE_NAME=quiethn-updater
UPDATER_BIN=${BIN_DIR}/${UPDATER_EXECUTABLE_NAME}

.PHONY: all
all: build

.PHONY: build
build: ${BIN} ${UPDATER_BIN}

.PHONY: clean
clean:
	go clean
	rm -rf ./bin

.PHONY: go-prepare
go-prepare:
	go mod tidy
	go mod download
	go test ./...

#
# quiethn
#

${BIN}: go-prepare
${BIN}: ${GO_MAIN} ${GO_SOURCES}
	CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -o $@ $<

.PHONY: run 
run: ${BIN}
	if [ ! -x $< ]; then chmod +x $<; fi
	HN_SERVER_PORT=8080 $<

#
# quiethn-updater
#

${UPDATER_BIN}: go-prepare
${UPDATER_BIN}: ${UPDATER_GO_MAIN} ${UPDATER_GO_SOURCES}
	CGO_ENABLED=0 GOOS=${UPDATER_TARGET_OS} GOARCH=${UPDATER_TARGET_ARCH} go build -o $@ $<

.PHONY: run-updater
run-updater: ${UPDATER_BIN}
	if [ ! -x $< ]; then chmod +x $<; fi
	HN_SERVER_PORT=8080 $<

#
# DOCKER
#

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
