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
RUN /usr/local/go/bin/go get github.com/golang/dep/cmd/dep && \
    go env


#CMD ["$GOPATH/src/github.com/cypherium/cypherscan-server/app"]
