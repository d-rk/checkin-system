FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go generate ./...

RUN go build -o checkin-system ./cmd/...

FROM alpine:latest
COPY --from=builder /app/checkin-system /checkin-system
COPY --from=builder /app/db /db

RUN touch .env

EXPOSE 8080

ENTRYPOINT ["/checkin-system"]
