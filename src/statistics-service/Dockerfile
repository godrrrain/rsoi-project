FROM golang:1.24.0 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o statistics-service ./src/statistics-service

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/statistics-service .

EXPOSE 8040

CMD ["./statistics-service"]
