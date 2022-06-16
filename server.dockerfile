FROM golang:1.18.3-alpine3.16 as build
WORKDIR /App
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /serverApp ./gin-server/main.go 


FROM scratch
COPY --from=build /serverApp /serverApp
ENTRYPOINT ["/serverApp"]

# cmd: docker build -f server.dockerfile .  -t server_service
# docker run --rm -p 1231:1231 -d server_service