FROM golang:alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o proxy

CMD ["./proxy"]
