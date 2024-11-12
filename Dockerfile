FROM golang:1.23-alpine AS builder

WORKDIR /build

ADD go.mod .

RUN go mod download

COPY . .

RUN go build -o playlist .

FROM alpine:3.20

WORKDIR /build

COPY --from=builder /build/playlist /build/playlist

RUN apk update \
    && apk -U upgrade \
    && apk --update add --no-cache dumb-init ca-certificates openssl \
    && update-ca-certificates \
    && chmod +x /build/playlist

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./playlist"]