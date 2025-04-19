FROM golang:1.23-alpine AS builder

WORKDIR /app/
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY ./backend .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /app/main .

COPY --from=builder /app/cookmeet-ai-b1a34baf28a6.json .
RUN mkdir -p cuisine_images user_images
EXPOSE 8080
CMD ["./main"]