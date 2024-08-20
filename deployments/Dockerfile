FROM golang:latest

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ../ ./
RUN go build -o main ./cmd/app

COPY ./configs/config.yaml ./config.yaml

EXPOSE 8080

CMD ["./main"]