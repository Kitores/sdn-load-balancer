FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN GOOS=linux go build -o main service1.go

FROM alpine
COPY --from=builder /app/main /app/main
ENTRYPOINT ["/app/main"]
EXPOSE 8080