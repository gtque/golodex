# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

RUN apk add build-base

COPY go.mod ./
COPY go.sum ./
COPY internal/ ./internal/
COPY pages/ ./pages/
COPY main.go ./
RUN go mod download

COPY build.sh build.sh
RUN chmod +x build.sh
RUN /app/build.sh

EXPOSE 8090

CMD [ "./golodex_api" ]