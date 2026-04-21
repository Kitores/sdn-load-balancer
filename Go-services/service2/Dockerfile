FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main service2.go

FROM alpine
COPY --from=builder /app/main /app/main
ENTRYPOINT ["/app/main"]
EXPOSE 8080