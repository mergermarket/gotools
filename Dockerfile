FROM golang:1.9

RUN go get -u github.com/alecthomas/gometalinter
RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/securego/gosec
RUN go get -u github.com/stretchr/testify/assert
RUN go get -u github.com/kyoh86/richgo
RUN gometalinter --install
WORKDIR /go/src/github.com/mergermarket/gotools
ADD . /go/src/github.com/mergermarket/gotools
CMD ./build.sh
