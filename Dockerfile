FROM golang:latest

RUN mkdir -p /usr/src/fibonacci
WORKDIR /usr/src/fibonacci

COPY ./ /usr/src/fibonacci
RUN go build -v -o ./bin/fibonacci ./cmd/fibonacci

EXPOSE 8080
EXPOSE 50052

ENV TZ Europe/Moscow

CMD ["./bin/fibonacci"]
