FROM ubuntu:16.04
RUN apt-get update  \
    && apt-get install -y gcc cmake libssl-dev openssl libgmp-dev bzip2 m4 build-essential git curl gcc libc-dev wget texinfo

#RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN  wget https://storage.googleapis.com/golang/go1.10.3.linux-amd64.tar.gz && \
     tar -C /usr/local -xzf go1.10.3.linux-amd64.tar.gz && \
     rm go1.10.3.linux-amd64.tar.gz && \
      echo 'export GOROOT=/usr/local/go' >> /etc/profile && \
      echo 'export GOPATH=$HOME/work' >> /etc/profile && \
      echo 'export GOBIN=$GOPATH/bin' >> /etc/profile && \
      echo 'export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH' >> ~/.bashrc && \
      /bin/bash -c "source ~/.bashrc"
RUN /usr/local/go/bin/go get github.com/golang/dep/cmd/dep
RUN wget https://ftp.gnu.org/gnu/gmp/gmp-6.1.2.tar.bz2 && \
    tar -xjf gmp-6.1.2.tar.bz2 && \
    cd gmp-6.1.2 && \
    ./configure --prefix=/usr --enable-cxx --disable-static --docdir=/usr/share/doc/gmp-6.1.2 && \
     make && \
     make check && \
     make install && \
     cp -rf /usr/lib/libgmp* /usr/local/lib/
RUN mkdir /root/go/src/github.com/cypherium -p && \
    cd /root/go/src/github.com/cypherium && \
    git clone https://258b8e7dc26fbd64e90e96d2c4290f3f81db1e9d@github.com/cypherium/cypherBFT.git --branch dTN-0.3

RUN mkdir /root/go/src/github.com/cypherium -p && \
    cd /root/go/src/github.com/cypherium && \
    git clone https://258b8e7dc26fbd64e90e96d2c4290f3f81db1e9d@github.com/cypherium/cypherscan-server.git && \
    cd cypherscan-server/cmd/main/ && \
    /root/go/bin/dep ensure && \
    /usr/local/go/bin/go build -o scan ./*
    /usr/local/go/bin/go build -o syncer ./*
    /usr/local/go/bin/go build -o query ./*

CMD ["/root/go/src/github.com/cypherium/cypherscan-server/cmd/main/scan"]
