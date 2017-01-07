FROM golang:1.7.3
COPY ./ /go/src/github.com/draganm/web-interceptor
WORKDIR /go/src/github.com/draganm/web-interceptor
RUN go build .
CMD ["/go/bin/web-interceptor"]
