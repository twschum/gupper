FROM golang:latest
LABEL maintainer="Tim Schumacher <twschum@gmail.com>"

WORKDIR /go/src/github.com/twschum/gupper
COPY cmd cmd
COPY pkg pkg

RUN go build -o server cmd/server/server.go

VOLUME /var/packages
EXPOSE 8080

CMD ["./server","-pkgdir","/var/packages"]
