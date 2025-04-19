FROM golang:1.23-alpine

WORKDIR /app/
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY ./backend .
# ビルド時のメモリ使用量を抑えるフラグを追加
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /app/main .

COPY --from=builder /app/cookmeet-ai-b1a34baf28a6.json .

# GCのガベージコレクション設定を調整
ENV GOGC=20
# 最大プロセス数を制限
ENV GOMAXPROCS=1

EXPOSE 8080
CMD ["go", "run", "main.go"]