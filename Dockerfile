FROM golang:1.24.5-alpine3.22

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 go build -o main ./cmd/server

EXPOSE 50051

CMD ["./main"]
