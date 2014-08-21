FROM mischief/docker-golang
ENV HOME /root
RUN go get github.com/dfjones/solar/solar-server
EXPOSE 3000

CMD $GOPATH/bin/solar-server