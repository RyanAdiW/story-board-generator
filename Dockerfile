FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/api ./cmd/api
RUN go build -o /bin/worker ./cmd/worker

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /bin/api /usr/local/bin/api
COPY --from=builder /bin/worker /usr/local/bin/worker

CMD ["api"]
