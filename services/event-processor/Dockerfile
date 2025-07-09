# Dockerfile
FROM golang:1.21.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o event-processor

# Final image
FROM gcr.io/distroless/base-debian11:latest

WORKDIR /app
COPY --from=builder /app/event-processor .

ENTRYPOINT ["/app/event-processor"]
