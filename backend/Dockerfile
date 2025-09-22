FROM golang:1.25-alpine AS builder

WORKDIR /app

ARG VERSION
ARG GIT_COMMIT

RUN apk add --no-cache make

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN make generate

RUN make build VERSION=$VERSION COMMIT_SHA=$GIT_COMMIT

FROM alpine:latest

RUN touch .env
RUN apk add --no-cache bash sshpass openssh

RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh

COPY --from=builder /app/checkin-system /checkin-system
COPY --from=builder /app/db /db

EXPOSE 8080

ENTRYPOINT ["/checkin-system"]
