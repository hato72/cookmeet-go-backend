FROM golang:1.23-alpine AS builder

WORKDIR /app/
COPY backend/go.mod backend/go.sum ./
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

ENV TEST_PORT=8081
ENV TEST_POSTGRES_USER=hato
ENV TEST_POSTGRES_PW=hato72
ENV TEST_POSTGRES_DB=hato_test
ENV TEST_POSTGRES_PORT=5433
ENV TEST_POSTGRES_HOST=postgres
ENV TEST_GO_ENV=test
ENV TEST_SECRET=test_secret
ENV TEST_API_DOMAIN=localhost
ENV TEST_FE_URL=http://localhost:3000

EXPOSE 8081
CMD ["./main"]