FROM golang:1.18.3-alpine3.16 as build
WORKDIR /App
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /orderApp ./order-service/main.go 


FROM scratch
COPY --from=build /orderApp /orderApp
ENTRYPOINT ["/orderApp"]

# cmd: docker build -f order.dockerfile .  -t order_service
# docker run --rm -p 50054:50054 -d order_service