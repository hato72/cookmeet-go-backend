FROM golang:1.22-alpine

WORKDIR /app/

# ローカルの場合

COPY backend/go.sum . 
COPY backend/go.mod .


# RUN apk upgrade --update && apk --no-cache add git gcc musl-dev

COPY ./backend .

RUN go build -tags netgo -ldflags '-s -w' -o app

EXPOSE 8081

CMD ["go", "run", "main.go"]



#本番環境の場合


# COPY go.mod ./

# RUN go mod download

# COPY . .

# #RUN go build -o main .

# EXPOSE 8080

# CMD ["go","run","./src/main.go"]


