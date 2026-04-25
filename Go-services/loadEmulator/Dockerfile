FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main emulator.go

FROM alpine
RUN apk add --no-cache bash
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
EXPOSE 8080