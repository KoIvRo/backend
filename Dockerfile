FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.23

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /server /server

USER app

EXPOSE 8000

ENTRYPOINT [ "/server" ]
