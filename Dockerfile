# Build Stage
FROM golang:1.23.2-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/rest/main.go

# Run Stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT [ "/app/main" ]
