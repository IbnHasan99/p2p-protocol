FROM golang:1.24.4

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o app ./cmd

ENTRYPOINT ["./app"]