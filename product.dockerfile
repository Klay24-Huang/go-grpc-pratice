FROM golang:1.18.3-alpine3.16 as build
WORKDIR /App
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /productApp ./product-service/main.go 


FROM scratch
COPY --from=build /productApp /productApp
ENTRYPOINT ["/productApp"]

# cmd: docker build -f product.dockerfile .  -t product_service