#FROM golang:alpine
FROM ubuntu:16.04
RUN apk update && \
    apk upgrade && \
    apk add git curl gcc libc-dev
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN apt-get update  \
    && apt-get install -y libssl-dev openssl libgmp-dev bzip2 m4

RUN && \
    wget https://ftp.gnu.org/gnu/gmp/gmp-6.1.2.tar.bz2 && \
    tar -xjf gmp-6.1.2.tar.bz2 && \
    cd gmp-6.1.2 && \
    ./configure --prefix=/usr --enable-cxx --disable-static --docdir=/usr/share/doc/gmp-6.1.2 && \
    sudo make && \
    sudo make html && \
    sudo make check && \
    sudo make install && \
    sudo make install-html && \
    sudo cp -rf /usr/lib/libgmp* /usr/local/lib/
RUN mkdir /go/src/github.com/cypherium -p && \
    cd /go/src/github.com/cypherium && \
    git clone https://258b8e7dc26fbd64e90e96d2c4290f3f81db1e9d@github.com/cypherium/cypherBFT.git --branch dTN-0.3

RUN mkdir /go/src/gitlab.com/ron-liu -p && \
    cd /go/src/gitlab.com/ron-liu && \
    git clone https://258b8e7dc26fbd64e90e96d2c4290f3f81db1e9d@github.com/cypherium/cypherscan-server.git && \
    cd cypherscan-server/cmd/main/ && \
    dep ensure && \
    go build -o app ./*

EXPOSE 8000

CMD ["/go/src/gitlab.com/ron-liu/cypherscan-server/app"]
