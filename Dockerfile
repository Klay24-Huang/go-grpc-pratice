# Golang Server's Dockerfile
FROM golang:1.14.6-alpine

RUN apk add --no-cache git

RUN go get github.com/lib/pq