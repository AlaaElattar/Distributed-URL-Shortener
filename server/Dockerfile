FROM golang:1.22

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o url-shortener

EXPOSE 8080

CMD ["./url-shortener"]