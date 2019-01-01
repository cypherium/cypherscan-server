FROM golang:alpine

WORKDIR /src
COPY ./

RUN go build  -o cypher-server

EXPOSE 8081
CMD ["./eb-go-app"]