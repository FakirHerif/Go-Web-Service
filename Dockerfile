FROM golang:1.21.4

WORKDIR /Application

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o go-web-service .

EXPOSE 8080

CMD ["/Application/go-web-service"]