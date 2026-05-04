FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main service1.go

FROM alpine
RUN apk add --no-cache bash iproute2 curl
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
EXPOSE 8080
