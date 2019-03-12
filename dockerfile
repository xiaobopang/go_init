FROM golang:1.11

WORKDIR $GOPATH/src/go_init
COPY . $GOPATH/src/go_init

RUN go build .

EXPOSE 7777
ENTRYPOINT ["./go_init"]