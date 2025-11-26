FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app 
RUN apk add --no-cache make
COPY . .
ARG TARGETARCH
RUN make TARGET_OS=linux TARGET_ARCH=${TARGETARCH}

FROM alpine:3.22

RUN apk add --no-cache curl

COPY --from=builder /app/bin/quiethn /bin/quiethn 

HEALTHCHECK --interval=15s --timeout=10s --start-period=5s --retries=3 CMD [ "/bin/sh", "-c", "curl -f 'http://localhost' || exit 1" ]

CMD [ "/bin/quiethn" ]
