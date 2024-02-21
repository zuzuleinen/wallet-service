FROM golang:1.22

RUN apt-get update && \
    apt-get install -y sqlite3 libsqlite3-dev

WORKDIR /app
COPY . .
RUN go build -o main

EXPOSE 8080

CMD ["./main"]
