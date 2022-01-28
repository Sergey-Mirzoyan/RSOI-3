FROM golang:latest as builder
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GO111MODULE=on

WORKDIR /app

ARG SERVICE
COPY ./src/$SERVICE ./src/$SERVICE
COPY ./src/database ./src/database
COPY ./src/utility ./src/utility
RUN go mod init lab2
RUN go mod tidy

WORKDIR ./src/$SERVICE
RUN go build -o main main.go

CMD ./main -docker
