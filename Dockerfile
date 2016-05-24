FROM ubuntu

RUN apt-get update
RUN apt-get install -y git mercurial curl

RUN curl https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz | \
    tar -C /usr/local -xz

ENV GOPATH $HOME/go
ENV PATH $PATH:/usr/local/go/bin:$HOME/go/bin

ADD . $GOPATH/src/github.com/davemt/dumbdumb
WORKDIR $GOPATH/src/github.com/davemt/dumbdumb

RUN go get ./...
RUN cd dumbdumb && go install

EXPOSE 25

CMD dumbdumb
