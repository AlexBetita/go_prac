FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-server ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth-server .

EXPOSE 8080
CMD ["./auth-server"]
