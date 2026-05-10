FROM golang:1.24.0

COPY ./ /app

RUN export GOPATH=/app

WORKDIR /app

RUN go mod tidy

RUN go build -o library ./src/library-service

ENTRYPOINT [ "./library" ]