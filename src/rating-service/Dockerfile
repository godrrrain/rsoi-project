FROM golang:1.24.0 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rating-service ./src/rating-service

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/rating-service .

EXPOSE 8090

CMD ["./rating-service"]