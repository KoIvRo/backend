FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate && \
    CGO_ENABLED=0 GOOS=linux go build -o /seed ./cmd/seed

FROM alpine:3.23

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /server /server
COPY --from=builder /migrate /migrate
COPY --from=builder /seed /seed

USER app

EXPOSE 8000

ENTRYPOINT [ "/server" ]
