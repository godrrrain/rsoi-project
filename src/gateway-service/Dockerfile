FROM golang:1.24.0 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway-service ./src/gateway-service

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/gateway-service .

EXPOSE 8090

CMD ["./gateway-service"]