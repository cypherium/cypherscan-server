FROM golang:alpine

RUN apk update && \
    apk upgrade && \
    apk add git curl gcc libc-dev

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh


RUN mkdir /go/src/github.com/cypherium -p && \
    cd /go/src/github.com/cypherium && \
    git clone https://258b8e7dc26fbd64e90e96d2c4290f3f81db1e9d@github.com/cypherium/CypherTestNet.git --branch reconfigTestNet

RUN mkdir /go/src/gitlab.com/ron-liu -p && \
    cd /go/src/gitlab.com/ron-liu && \
    git clone https://258b8e7dc26fbd64e90e96d2c4290f3f81db1e9d@github.com/cypherium/cypherscan-server.git && \
    cd cypherscan-server && \
    dep ensure && \
    go build -o app cmd/main/*

CMD ["/go/src/gitlab.com/ron-liu/cypherscan-server/app"]