# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/src/ -v ./...

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/src/. /exchangeapi
COPY --from=builder /go/src/app/resources/application.yaml /resources/application.yaml
EXPOSE 8080
CMD /exchangeapi/exchange-web-service
