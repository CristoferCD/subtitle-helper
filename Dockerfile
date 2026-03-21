FROM golang:1.25.6 AS builder

WORKDIR /app

COPY go.mod go.sum *.go .
RUN go mod download
RUN CGO_ENABLED=0 \
    go build -o main *.go

FROM alpine:3.23

RUN apk upgrade -U \ 
    && apk add ca-certificates ffmpeg \
    && rm -rf /var/cache/*

COPY --from=builder /app/main /app/main

CMD ["/app/main"]