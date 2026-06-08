FROM golang:1.24.9-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o account-service .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /build/account-service .

USER app

EXPOSE 8080 9090

ENTRYPOINT ["./account-service"]
