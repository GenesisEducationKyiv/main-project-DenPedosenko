# Stage 1: Build stage
FROM golang:latest AS builder
LABEL maintainer="logs-consumer"
WORKDIR /logs-consumer
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main main.go

# Stage 2: Final stage
FROM golang:latest
WORKDIR /logs-consumer
COPY --from=builder /logs-consumer/main .
CMD ["./main"]
EXPOSE 8080
