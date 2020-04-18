FROM ubuntu:16.04

ARG EXECUTION_TIMEOUT=default_value
ENV EXECUTION_TIMEOUT=$EXECUTION_TIMEOUT
ARG NODES_URLS=default_value
ENV NODES_URLS=$NODES_URLS
ARG DYNAMODB_REGION=default_value
ENV DYNAMODB_REGION=$DYNAMODB_REGION
ARG AWS_ACCESS_KEY_ID=default_value
ENV AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY=default_value
ENV AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
ARG REGION=default_value
ENV REGION=$REGION
ARG RECENT_TTL_DURATION_IN_SECONDS=default_value
ENV RECENT_TTL_DURATION_IN_SECONDS=$RECENT_TTL_DURATION_IN_SECONDS

RUN apt-get update  \
    && apt-get install -y gcc cmake libssl-dev openssl libgmp-dev bzip2 m4 build-essential git curl gcc libc-dev wget texinfo
RUN mkdir /root/go/src/github.com/cypherium -p && \
    cd /root/go/src/github.com/cypherium && \
    git clone https://258b8e7dc26fbd64e90e96d2c4290f3f81db1e9d@github.com/cypherium/cypherscan-server.git --branch scan

#RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN  wget https://storage.googleapis.com/golang/go1.10.3.linux-amd64.tar.gz && \
     tar -C /usr/local -xzf go1.10.3.linux-amd64.tar.gz && \
     rm go1.10.3.linux-amd64.tar.gz
#      echo 'export GOROOT=/usr/local/go' >> /etc/profile && \
#      echo 'export GOPATH=$HOME/work' >> /etc/profile && \
#      echo 'export GOBIN=$GOPATH/bin' >> /etc/profile && \
#      echo 'export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH' >> /etc/profile && \
#      /bin/bash -c "source /etc/profile"

      echo 'export EXECUTION_TIMEOUT=0' >> /etc/profile && \
      echo 'export NODES_URLS=ws://40.117.112.213:8546' >> /etc/profile && \
      echo 'export DYNAMODB_REGION=us-east-2' >> /etc/profile && \
      echo 'export AWS_ACCESS_KEY_ID=AKIAJYWTBXV3Z2HWLE3Q' >> /etc/profile && \
      echo 'export AWS_SECRET_ACCESS_KEY=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' >> /etc/profile && \
      echo 'export RECENT_TTL_DURATION_IN_SECONDS=36000000' >> /etc/profile && \
      /bin/bash -c "source /etc/profile"
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
    cd cypherscan-server/cmd/main/ && \
    /root/go/bin/dep ensure && \
    /usr/local/go/bin/go build -o scan ./*

CMD ["/root/go/src/github.com/cypherium/cypherscan-server/cmd/main/scan"]
