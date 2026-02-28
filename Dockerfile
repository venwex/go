# ---------- STAGE 1: BUILD ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# install git (нужно для go mod download иногда)
RUN apk add --no-cache git

# сначала зависимости (Docker cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# потом код
COPY . .

# собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api


# ---------- STAGE 2: RUNTIME ----------
FROM alpine:3.19

WORKDIR /app

# копируем только бинарник
COPY --from=builder /app/app .
COPY --from=builder /app/database ./database

EXPOSE 8080

CMD ["./app"]