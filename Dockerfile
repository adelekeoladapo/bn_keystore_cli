FROM golang:1.18
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go test ./...
RUN go build -o main ./cmd
CMD ["/app/main"]