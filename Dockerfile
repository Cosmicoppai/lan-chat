FROM golang:1.17-alpine

EXPOSE 80
EXPOSE 9000

RUN mkdir -p lan-chat/movie

WORKDIR lan-chat

COPY . .

RUN go build -o start-chat

CMD ["start-chat"]
