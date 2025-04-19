FROM golang:1.24-alpine

WORKDIR /app/

COPY backend/go.mod .

RUN go mod download

COPY ./backend .

#RUN go build -o main .

EXPOSE 8080

CMD ["go", "run", "main.go"]

