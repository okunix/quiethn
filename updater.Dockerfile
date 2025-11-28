FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app
RUN apk add --no-cache make
COPY . .
ARG TARGETARCH
RUN make TARGET_OS=linux TARGET_ARCH=${TARGETARCH}

FROM alpine:3.22

COPY --from=builder /app/bin/quiethn-updater /bin/quiethn-updater
ENTRYPOINT [ "/bin/quiethn-updater" ]
