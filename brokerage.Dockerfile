# syntax=docker/dockerfile:1

FROM golang
WORKDIR /brokerage/build
COPY . .
RUN go mod download
RUN go build -o /brokerage-app