FROM golang:1.18.3-alpine3.16 as build
WORKDIR /App
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /userApp ./user-service/main.go 


FROM scratch
COPY --from=build /userApp /userApp
ENTRYPOINT ["/userApp"]

# cmd: docker build -f user.dockerfile .  -t user_service
# docker run --rm -p 50052:50052 -d user_service