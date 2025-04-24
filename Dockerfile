FROM golang:1.23-alpine AS builder

WORKDIR /app/
COPY ./backend/go.mod ./backend/go.sum ./
COPY ./backend/.env.test ./
RUN go mod download
COPY ./backend .
# ビルド時のメモリ使用量を抑えるフラグを追加
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /root/
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/main .

# GCのガベージコレクション設定を調整
ENV GOGC=20
# 最大プロセス数を制限
ENV GOMAXPROCS=1

EXPOSE 8081
CMD ["./main"]