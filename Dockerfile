FROM golang:1.26.2 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/posts ./cmd/app

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/posts /app/posts

CMD ["./posts"]
