FROM mischief/docker-golang
ENV HOME /root
RUN go get github.com/dfjones/solar/server

CMD go run $GOPATH/src/github.com/dfjones/solar/server/server.go