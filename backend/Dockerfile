FROM golang:1.19-alpine

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /checkin-system

RUN touch .env

EXPOSE 8080

CMD [ "/checkin-system" ]